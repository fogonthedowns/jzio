import aws_cdk as cdk
from aws_cdk import pipelines, aws_iam as iam
from constructs import Construct

from stacks.stage_site import SiteStage

_REPO = "fogonthedowns/jzio"
_BRANCH = "master"

# Files/dirs that live in the repo root but must NOT be synced to S3.
_EXCLUDES = (
    '--exclude ".git/*" '
    '--exclude "infra/*" '
    '--exclude ".cursor/*" '
    '--exclude ".gitignore" '
    '--exclude ".arcconfig" '
    '--exclude "node_modules/*"'
)


class PipelineStack(cdk.Stack):
    def __init__(self, scope: Construct, id: str, *, connection_arn: str, **kwargs):
        super().__init__(scope, id, **kwargs)

        source = pipelines.CodePipelineSource.connection(
            _REPO, _BRANCH,
            connection_arn=connection_arn,
        )

        pipeline = pipelines.CodePipeline(
            self, "Pipeline",
            pipeline_name="jzio-pipeline",
            cross_account_keys=False,
            synth=pipelines.ShellStep(
                "Synth",
                input=source,
                commands=[
                    "npm install -g aws-cdk",
                    "cd infra && python3 -m venv .venv && .venv/bin/pip install -q -r requirements.txt && cdk synth",
                ],
                primary_output_directory="infra/cdk.out",
            ),
        )

        stage = SiteStage(self, "Prod")

        # After SiteStack deploys: sync static files and bust the CDN cache.
        sync_step = pipelines.CodeBuildStep(
            "SyncSite",
            input=source,
            env_from_cfn_outputs={
                "BUCKET_NAME": stage.site_stack.bucket_name_cfn_output,
                "DIST_ID": stage.site_stack.distribution_id_cfn_output,
            },
            commands=[
                f'aws s3 sync . "s3://$BUCKET_NAME/" --region us-west-2 --delete {_EXCLUDES}',
                'aws cloudfront create-invalidation --region us-west-2 --distribution-id "$DIST_ID" --paths "/*"',
            ],
            role_policy_statements=[
                iam.PolicyStatement(
                    actions=["s3:ListBucket", "s3:GetBucketLocation"],
                    resources=[stage.site_stack.site_bucket.bucket_arn],
                ),
                iam.PolicyStatement(
                    actions=["s3:GetObject", "s3:PutObject", "s3:DeleteObject"],
                    resources=[stage.site_stack.site_bucket.arn_for_objects("*")],
                ),
                iam.PolicyStatement(
                    actions=["cloudfront:CreateInvalidation"],
                    resources=["*"],
                ),
            ],
        )

        pipeline.add_stage(stage, post=[sync_step])
