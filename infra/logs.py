import boto3
import sys

BUCKET_NAME = 'jzio-logs'
PREFIX = 'cf/'

def format_log_entry(date, time, method, path, status, result, bytes_sent, ip):
    """Format a single log entry for display."""
    timestamp = f"{date} {time}"
    return f"{timestamp} | {status} {result:8} | {method:4} {path:40} | {bytes_sent:>6}B | {ip}"

def query_logs_with_s3_select(where_clause=""):
    """Query CloudFront logs using S3 Select."""
    s3 = boto3.client('s3')

    # CloudFront logs use tab-separated values with a custom header
    sql = f"""
    SELECT _1 as date, _2 as time, _6 as method, _7 as path,
           _8 as status, _15 as result, _3 as bytes, _5 as ip
    FROM s3object
    WHERE _1 != '#date' {where_clause}
    """

    paginator = s3.get_paginator('list_objects_v2')
    page_iterator = paginator.paginate(Bucket=BUCKET_NAME, Prefix=PREFIX)

    total_requests = 0
    status_counts = {}
    result_counts = {}

    print(f"📜 CloudFront Logs from {BUCKET_NAME}/{PREFIX}\n")
    print(f"{'Timestamp':<19} | {'Status':<10} | {'Method Path':<45} | {'Bytes':>6} | IP")
    print("─" * 120)

    for page in page_iterator:
        if 'Contents' not in page:
            continue

        for obj in page['Contents']:
            s3_key = obj['Key']
            if s3_key.endswith('/'):
                continue

            try:
                # Use S3 Select to query the gzipped log file
                response = s3.select_object_content(
                    Bucket=BUCKET_NAME,
                    Key=s3_key,
                    ExpressionType='SQL',
                    Expression=sql,
                    InputSerialization={
                        'CSV': {
                            'FileHeaderInfo': 'NONE',
                            'Comments': '#',
                            'AllowQuotedRecordDelimiter': False,
                            'RecordDelimiter': '\n',
                            'FieldDelimiter': '\t',
                        },
                        'CompressionType': 'GZIP',
                    },
                    OutputSerialization={'CSV': {}},
                )

                # Process the streamed response
                for event in response['Payload']:
                    if 'Records' in event:
                        payload = event['Records']['Payload'].decode('utf-8')
                        for line in payload.strip().split('\n'):
                            if line:
                                parts = line.split(',')
                                if len(parts) >= 8:
                                    date, time, method, path, status, result, bytes_sent, ip = parts[:8]
                                    # Remove quotes if present
                                    path = path.strip('"')
                                    print(format_log_entry(date, time, method, path, status, result, bytes_sent, ip))
                                    total_requests += 1
                                    status_counts[status] = status_counts.get(status, 0) + 1
                                    result_counts[result] = result_counts.get(result, 0) + 1

            except Exception as e:
                print(f"❌ Error querying {s3_key}: {e}")

    # Summary stats
    print("\n" + "─" * 120)
    print(f"\n📊 Summary:")
    print(f"   Total requests: {total_requests}")
    print(f"   Status codes: {dict(sorted(status_counts.items()))}")
    print(f"   Results: {dict(sorted(result_counts.items()))}")

def stream_s3_logs(status_filter=None, result_filter=None):
    """Query logs with optional filters.

    Args:
        status_filter: Status code to filter by (e.g., '200', '404')
        result_filter: CloudFront result type to filter by (e.g., 'Error', 'Hit')
    """
    where_parts = []

    if status_filter:
        where_parts.append(f"_8 = '{status_filter}'")
    if result_filter:
        where_parts.append(f"_15 = '{result_filter}'")

    where_clause = " AND ".join(where_parts)
    if where_clause:
        where_clause = "AND " + where_clause

    query_logs_with_s3_select(where_clause)

if __name__ == '__main__':
    # Optional command-line filters:
    # python logs.py              # All logs
    # python logs.py 200          # Only 200 status codes
    # python logs.py 404          # Only 404s
    # python logs.py 200 Hit      # 200 status AND Hit result
    status_filter = sys.argv[1] if len(sys.argv) > 1 else None
    result_filter = sys.argv[2] if len(sys.argv) > 2 else None
    stream_s3_logs(status_filter, result_filter)
