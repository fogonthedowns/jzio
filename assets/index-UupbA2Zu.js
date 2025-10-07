(function(){const r=document.createElement("link").relList;if(r&&r.supports&&r.supports("modulepreload"))return;for(const t of document.querySelectorAll('link[rel="modulepreload"]'))i(t);new MutationObserver(t=>{for(const o of t)if(o.type==="childList")for(const s of o.addedNodes)s.tagName==="LINK"&&s.rel==="modulepreload"&&i(s)}).observe(document,{childList:!0,subtree:!0});function n(t){const o={};return t.integrity&&(o.integrity=t.integrity),t.referrerPolicy&&(o.referrerPolicy=t.referrerPolicy),t.crossOrigin==="use-credentials"?o.credentials="include":t.crossOrigin==="anonymous"?o.credentials="omit":o.credentials="same-origin",o}function i(t){if(t.ep)return;t.ep=!0;const o=n(t);fetch(t.href,o)}})();const K="https://api.themoviedb.org/3",D="eyJhbGciOiJIUzI1NiJ9.eyJhdWQiOiI5NjUzODlhYzU2YWUzZjgzMGYxMTg5NzgxNDBiOWZjMSIsIm5iZiI6MTc1ODA1NTAxMS44NzksInN1YiI6IjY4YzljYTYzOWU1NzQ4OTlhZWRjZDU0ZiIsInNjb3BlcyI6WyJhcGlfcmVhZCJdLCJ2ZXJzaW9uIjoxfQ.vVSwVR_SreI6fbeiOrsCkekvSpVS-RrhsaU8R7dklK4";async function v(e,r={}){const n=new URL(`${K}${e}`);Object.entries(r).forEach(([t,o])=>{n.searchParams.append(t,o)});const i=await fetch(n.toString(),{headers:{Authorization:`Bearer ${D}`,accept:"application/json"}});if(!i.ok)throw new Error(`TMDb API error: ${i.status} ${i.statusText}`);return i.json()}async function F(e){return v("/search/movie",{query:e})}async function y(e){return v(`/movie/${e}`)}async function G(e){return v(`/movie/${e}/similar`)}async function Q(e){return v(`/movie/${e}/recommendations`)}async function N(){return v("/movie/popular")}async function T(){return v("/movie/top_rated")}const k="watched",w="watchedMovies",R="userName",q="recentlyViewed";let p=JSON.parse(localStorage.getItem(k)||"[]"),u=JSON.parse(localStorage.getItem(w)||"[]"),g=JSON.parse(localStorage.getItem(q)||"[]");function X(){return p}function ee(e){p.includes(e)||(p.push(e),localStorage.setItem(k,JSON.stringify(p)))}function f(){return u}function te(e){const r=u.find(n=>n.id===e.id);r?(r.rating=e.rating||r.rating||0,r.comment=e.comment||r.comment||"",localStorage.setItem(w,JSON.stringify(u))):(u.push({id:e.id,title:e.title,release_date:e.release_date,poster_path:e.poster_path,rating:e.rating||0,comment:e.comment||""}),localStorage.setItem(w,JSON.stringify(u)))}function re(e,r){const n={n:e,m:r.map(i=>({i:i.id,t:i.title,y:i.release_date?i.release_date.split("-")[0]:"",p:i.poster_path,r:i.rating||0,c:i.comment||""}))};return btoa(JSON.stringify(n)).replace(/[+/=]/g,i=>i==="+"?"-":i==="/"?"_":"")}function ne(e){try{const r=e.replace(/[-_]/g,t=>t==="-"?"+":"/"),n=r+"=".repeat((4-r.length%4)%4),i=JSON.parse(atob(n));return{name:i.n,movies:i.m.map(t=>({id:t.i,title:t.t,release_date:t.y?`${t.y}-01-01`:"",poster_path:t.p,rating:t.r||0,comment:t.c||""}))}}catch{return null}}function ie(e){u=e}function O(e){p=p.filter(r=>r!==e),localStorage.setItem(k,JSON.stringify(p))}function A(e){u=u.filter(r=>r.id!==e),localStorage.setItem(w,JSON.stringify(u))}function H(){return localStorage.getItem(R)}function ae(e){localStorage.setItem(R,e)}function j(e){g=g.filter(r=>r.id!==e.id),g.unshift({id:e.id,title:e.title,release_date:e.release_date,poster_path:e.poster_path}),g=g.slice(0,5),localStorage.setItem(q,JSON.stringify(g))}function oe(){return g}const se=document.querySelector("#search-form"),B=document.querySelector("#search-input"),l=document.querySelector("#results"),z=document.querySelector("#watched-list"),W=document.querySelector("#recently-viewed-list"),de=document.querySelector("#share-button-container"),d=document.querySelector("#modal"),M=document.querySelector("#my-movies-link");let b="";se.addEventListener("submit",async e=>{e.preventDefault();const r=B.value.trim();if(r){l.innerHTML='<div style="text-align: center; padding: 2rem; color: rgba(0, 0, 0, 0.7);" class="loading">Searching for movies...</div>';try{const n=await F(r);n.results.length===0?l.innerHTML='<div style="text-align: center; padding: 2rem; color: rgba(0, 0, 0, 0.7);">No movies found. Try a different search term.</div>':C(n.results)}catch(n){console.error(n),l.innerHTML='<div style="text-align: center; padding: 2rem; color: rgba(255, 107, 107, 0.8);">⚠️ Something went wrong! Please try again.</div>'}}});function C(e,r=!1){const n=e.map(i=>`
      <div class="card" data-id="${i.id}">
        <img src="https://image.tmdb.org/t/p/w200${i.poster_path}" alt="${i.title}" />
        <div class="card-content">
          <h3>${i.title}</h3>
          <p>${i.release_date?i.release_date.split("-")[0]:"N/A"}</p>
        </div>
      </div>
    `).join("");r?l.innerHTML+=n:l.innerHTML=n}l.addEventListener("click",async e=>{const r=e.target.closest(".card");if(!r)return;const n=Number(r.dataset.id),i=await y(n);J(i)});function J(e){j(e),I(),!b&&(d.innerHTML=`
    <div id="modal-content">
      <h2>${e.title} (${e.release_date?.split("-")[0]||""})</h2>
      <img src="https://image.tmdb.org/t/p/w300${e.poster_path}" />
      <p>${e.overview}</p>
      <button id="add-watched">Add to Watched</button>
      <button id="close-modal">Close</button>
    </div>
  `,d.classList.remove("hidden"),document.querySelector("#close-modal").addEventListener("click",()=>d.classList.add("hidden")),document.querySelector("#add-watched").addEventListener("click",()=>{Y(e)}))}function Y(e){d.innerHTML=`
    <div class="dialog-content">
      <h3>Rate & Review "${e.title}"</h3>
      <div class="star-rating" id="star-rating">
        ${Array(5).fill(0).map((t,o)=>`<span class="star" data-rating="${o+1}">★</span>`).join("")}
      </div>
      <textarea class="dialog-textarea" id="review-comment" placeholder="What did you think about this movie? (optional)"></textarea>
      <div class="dialog-buttons">
        <button class="dialog-btn dialog-btn-secondary" id="cancel-review">Cancel</button>
        <button class="dialog-btn dialog-btn-primary" id="save-review">Save Review</button>
      </div>
    </div>
  `;let r=0;const n=d.querySelectorAll(".star"),i=d.querySelector("#review-comment");n.forEach((t,o)=>{t.addEventListener("mouseover",()=>{n.forEach((s,a)=>{a<=o?s.classList.add("active"):s.classList.remove("active")})}),t.addEventListener("click",()=>{r=o+1,n.forEach((s,a)=>{a<r?s.classList.add("active"):s.classList.remove("active")})})}),d.addEventListener("mouseleave",()=>{n.forEach((t,o)=>{o<r?t.classList.add("active"):t.classList.remove("active")})}),document.getElementById("cancel-review").addEventListener("click",()=>{d.classList.add("hidden")}),document.getElementById("save-review").addEventListener("click",()=>{const t={...e,rating:r,comment:i.value.trim()};ee(e.id),te(t),_(),d.classList.add("hidden"),h("Movie added to your watched list!","success"),!B.value.trim()&&!b&&setTimeout(()=>Z(),1e3)})}function _(){const e=f();if(b){le(e);return}if(e.length===0){z.innerHTML='<li style="text-align: center; color: rgba(0, 0, 0, 0.5); padding: 1rem;">No movies watched yet</li>';return}const n=e.slice(-2).map(t=>`
      <li class="watched-movie-item" data-movie-id="${t.id}" style="position: relative; cursor: pointer; transition: all 0.2s ease; margin-bottom: 0.5rem;">
        <div style="display: flex; gap: 0.75rem; padding: 0.5rem;">
          <img src="https://image.tmdb.org/t/p/w92${t.poster_path}"
               style="width: 50px; height: 75px; border-radius: 6px; object-fit: cover; flex-shrink: 0;" />
          <div style="flex: 1; min-width: 0;">
            <div style="font-weight: 600; font-size: 0.95rem; margin-bottom: 0.3rem; line-height: 1.2; color: black;">${t.title}</div>
            <div style="font-size: 0.8rem; color: rgba(0, 0, 0, 0.6); margin-bottom: 0.5rem;">${t.release_date?t.release_date.split("-")[0]:"N/A"}</div>

            ${t.rating>0?`
              <div style="font-size: 0.85rem; color: #ffd700; display: flex; align-items: center; gap: 0.5rem; margin-bottom: 0.5rem;">
                <span>${"★".repeat(t.rating)}${"☆".repeat(5-t.rating)}</span>
                <span style="color: rgba(0,0,0,0.8); font-size: 0.75rem; font-weight: 600;">${t.rating}/5</span>
              </div>
            `:""}

            ${t.comment?`
              <div style="font-size: 0.8rem; color: rgba(0, 0, 0, 0.8); line-height: 1.3; font-style: italic; background: rgba(0,0,0,0.05); padding: 0.4rem 0.6rem; border-radius: 6px; border-left: 2px solid rgb(255, 102, 0);">
                "${t.comment.length>60?t.comment.substring(0,60)+"...":t.comment}"
              </div>
            `:""}
          </div>
          <button class="delete-movie-btn" style="
            position: absolute;
            right: 0.5rem;
            top: 0.5rem;
            background: rgb(255, 102, 0);
            color: white;
            border: none;
            border-radius: 50%;
            width: 24px;
            height: 24px;
            cursor: pointer;
            font-size: 12px;
            display: none;
            align-items: center;
            justify-content: center;
            transition: all 0.2s ease;
            z-index: 10;
          ">✕</button>
        </div>
      </li>
    `).join("");e.length>0?M.textContent=`My Movies (${e.length})`:M.textContent="My Movies",z.innerHTML=n,de.innerHTML=`
    <button id="share-btn" style="
      width: 100%;
      margin-top: 1rem;
      padding: 0.75rem;
      background: white;
      color: black;
      border: none !important;
      outline: none !important;
      border-radius: 10px;
      font-weight: 600;
      cursor: pointer;
      transition: all 0.3s ease;
      box-shadow: 0 4px 15px rgba(0, 0, 0, 0.2);
    ">Share My List</button>
  `,document.getElementById("share-btn")?.addEventListener("click",V),document.querySelectorAll(".watched-movie-item").forEach(t=>{const o=t.querySelector(".delete-movie-btn");t.addEventListener("mouseenter",()=>{o.style.display="flex"}),t.addEventListener("mouseleave",()=>{o.style.display="none"}),t.addEventListener("click",async()=>{const s=Number(t.dataset.movieId),a=e.find(c=>c.id===s);if(a){const c=await y(s);U(c,a)}}),o.addEventListener("click",s=>{s.stopPropagation();const a=Number(t.dataset.movieId);ce(a)})})}function ce(e){O(e),A(e),_(),h("Movie removed from your list","success")}function P(){const e=f(),r=document.querySelector("main"),n=document.querySelector("header h1");n.innerHTML=`
    <button id="back-home" style="
      background: white;
      color: black;
      border: none !important;
      outline: none !important;
      padding: 0.5rem 1rem;
      border-radius: 25px;
      cursor: pointer;
      font-size: 0.9rem;
      margin-right: 1rem;
      transition: all 0.3s ease;
      box-shadow: 0 4px 15px rgba(0, 0, 0, 0.2);
      font-weight: 600;
    ">← Back</button>
    All My Watched Movies (${e.length})
    <button id="share-all-movies" style="
      background: white;
      color: black;
      border: none !important;
      outline: none !important;
      padding: 0.5rem 1rem;
      border-radius: 25px;
      cursor: pointer;
      font-size: 0.9rem;
      margin-left: 1rem;
      transition: all 0.3s ease;
      box-shadow: 0 4px 15px rgba(0, 0, 0, 0.2);
    ">Share My List</button>
  `;const i=document.querySelector("#search-form");i.style.display="none",r.innerHTML=`
    <div class="shared-movie-grid">
      ${e.map(t=>`
        <div class="shared-movie-card watched-movie-card" data-movie-id="${t.id}">
          <img src="https://image.tmdb.org/t/p/w500${t.poster_path}"
               alt="${t.title}"
               class="shared-movie-image" />
          <div class="shared-movie-info">
            <h3 class="shared-movie-title">${t.title}</h3>
            <p class="shared-movie-year">${t.release_date?t.release_date.split("-")[0]:"N/A"}</p>

            ${t.rating>0?`
              <div class="shared-movie-rating">
                <span class="shared-movie-stars">${"★".repeat(t.rating)}${"☆".repeat(5-t.rating)}</span>
                <span>${t.rating}/5</span>
              </div>
            `:""}

            ${t.comment?`
              <div class="shared-movie-comment">
                "${t.comment}"
              </div>
            `:""}

            <button class="delete-watched-btn" style="
              position: absolute;
              top: 10px;
              right: 10px;
              background: rgba(255, 107, 107, 0.9);
              color: white;
              border: none;
              border-radius: 50%;
              width: 30px;
              height: 30px;
              cursor: pointer;
              font-size: 14px;
              display: flex;
              align-items: center;
              justify-content: center;
              transition: all 0.2s ease;
              opacity: 0;
            " onmouseenter="this.style.opacity='1'"
               onmouseleave="this.style.opacity='0'">✕</button>
          </div>
        </div>
      `).join("")}
    </div>
  `,document.getElementById("back-home").addEventListener("click",()=>{window.location.reload()}),document.getElementById("share-all-movies").addEventListener("click",V),document.querySelectorAll(".watched-movie-card").forEach(t=>{t.addEventListener("click",async s=>{if(s.target.classList.contains("delete-watched-btn"))return;const a=Number(s.currentTarget.dataset.movieId),c=e.find(m=>m.id===a);if(c){const m=await y(a);U(m,c)}}),t.querySelector(".delete-watched-btn").addEventListener("click",s=>{s.stopPropagation();const a=Number(t.dataset.movieId),c=e.find(m=>m.id===a)?.title;confirm(`Remove "${c}" from your watched list?`)&&(O(a),A(a),h("Movie removed from your list","success"),P())}),t.addEventListener("mouseenter",()=>{const s=t.querySelector(".delete-watched-btn");s.style.opacity="1"}),t.addEventListener("mouseleave",()=>{const s=t.querySelector(".delete-watched-btn");s.style.opacity="0"})})}function I(){const e=oe();if(e.length===0){W.innerHTML='<li style="text-align: center; color: rgba(0, 0, 0, 0.5); font-size: 0.9rem; padding: 1rem;">No recent views</li>';return}W.innerHTML=e.map(r=>`
      <li class="recent-movie-item" data-movie-id="${r.id}" style="
        cursor: pointer;
        transition: all 0.2s ease;
        margin-bottom: 0.5rem;
        padding: 0.5rem;
        border-radius: 8px;
      ">
        <div style="display: flex; align-items: center; gap: 0.5rem;">
          <img src="https://image.tmdb.org/t/p/w92${r.poster_path}"
               style="width: 35px; height: 52px; border-radius: 4px; object-fit: cover; flex-shrink: 0;" />
          <div style="flex: 1; min-width: 0;">
            <div style="font-weight: 600; font-size: 0.85rem; line-height: 1.2; margin-bottom: 0.2rem; overflow: hidden; text-overflow: ellipsis; white-space: nowrap;">
              ${r.title}
            </div>
            <div style="font-size: 0.75rem; color: rgba(0, 0, 0, 0.6);">
              ${r.release_date?r.release_date.split("-")[0]:"N/A"}
            </div>
          </div>
        </div>
      </li>
    `).join(""),document.querySelectorAll(".recent-movie-item").forEach(r=>{r.addEventListener("click",async()=>{const n=Number(r.dataset.movieId),i=await y(n);J(i)})})}function U(e,r){d.innerHTML=`
    <div id="modal-content" style="max-width: 800px; width: 100%;">
      <div style="display: flex; gap: 2rem; margin-bottom: 2rem;">
        <img src="https://image.tmdb.org/t/p/w400${e.poster_path}"
             style="width: 300px; height: 450px; object-fit: cover; border-radius: 15px; box-shadow: 0 10px 30px rgba(0,0,0,0.5);" />
        <div style="flex: 1;">
          <h2 style="margin: 0 0 1rem 0; font-size: 2.2rem; line-height: 1.2;">${e.title}</h2>
          <p style="color: rgba(0,0,0,0.8); font-size: 1.1rem; margin-bottom: 1.5rem;">
            ${e.release_date?e.release_date.split("-")[0]:"N/A"} •
            ${e.runtime?`${e.runtime} min`:""} •
            ${e.genres?.map(n=>n.name).join(", ")||""}
          </p>

          ${r.rating>0?`
            <div style="margin-bottom: 2rem; text-align: center; background: rgba(0,0,0,0.05); padding: 1.5rem; border-radius: 15px;">
              <div style="color: #ffd700; font-size: 2.5rem; margin-bottom: 0.5rem;">
                ${"★".repeat(r.rating)}${"☆".repeat(5-r.rating)}
              </div>
              <div style="font-size: 1.2rem; font-weight: 600; color: rgba(0,0,0,0.9);">
                Your Rating: ${r.rating}/5
              </div>
            </div>
          `:""}

          ${r.comment?`
            <div style="margin-bottom: 2rem;">
              <h3 style="color: rgb(255, 102, 0); margin-bottom: 1rem; font-size: 1.2rem;">Your Review</h3>
              <div style="background: rgba(0,0,0,0.05); border-left: 3px solid rgb(255, 102, 0); padding: 1rem; border-radius: 0 10px 10px 0; font-style: italic; line-height: 1.6; color: rgba(0,0,0,0.9);">
                "${r.comment}"
              </div>
            </div>
          `:""}

          ${e.vote_average?`
            <div style="margin-bottom: 1.5rem; display: flex; align-items: center; gap: 1rem;">
              <div style="background: rgba(255,193,7,0.2); padding: 0.5rem 1rem; border-radius: 25px; display: flex; align-items: center; gap: 0.5rem;">
                <span style="color: #ffc107;">⭐</span>
                <span style="font-weight: 600;">${e.vote_average.toFixed(1)}/10 TMDb</span>
              </div>
              <span style="color: rgba(0,0,0,0.7); font-size: 0.9rem;">(${e.vote_count} votes)</span>
            </div>
          `:""}
        </div>
      </div>

      <div style="margin-bottom: 2rem;">
        <h3 style="color: rgb(255, 102, 0); margin-bottom: 1rem;">Overview</h3>
        <p style="line-height: 1.6; color: rgba(0,0,0,0.9);">${e.overview}</p>
      </div>

      ${e.production_companies&&e.production_companies.length>0?`
        <div style="margin-bottom: 2rem;">
          <h4 style="color: rgba(0,0,0,0.8); margin-bottom: 0.5rem; font-size: 0.9rem;">Production</h4>
          <p style="color: rgba(0,0,0,0.6); font-size: 0.9rem;">
            ${e.production_companies.map(n=>n.name).join(", ")}
          </p>
        </div>
      `:""}

      <div style="text-align: center; margin-top: 2rem;">
        <button id="close-expanded-modal" style="
          background: rgba(0, 0, 0, 0.1);
          border: 1px solid rgba(0, 0, 0, 0.2);
          color: black;
          padding: 0.875rem 2rem;
          border-radius: 50px;
          font-weight: 600;
          cursor: pointer;
          transition: all 0.3s ease;
          font-size: 1rem;
        ">Close</button>
      </div>
    </div>
  `,d.classList.remove("hidden"),document.querySelector("#close-expanded-modal").addEventListener("click",()=>d.classList.add("hidden"))}function le(e){const r=document.querySelector("main"),n=document.querySelector("header h1");n.innerHTML=`
    <button id="back-home" style="
      background: white;
      color: black;
      border: none !important;
      outline: none !important;
      padding: 0.5rem 1rem;
      border-radius: 25px;
      cursor: pointer;
      font-size: 0.9rem;
      margin-right: 1rem;
      transition: all 0.3s ease;
      box-shadow: 0 4px 15px rgba(0, 0, 0, 0.2);
      font-weight: 600;
    ">← Back Home</button>
    🎬 ${b}'s Movies
  `;const i=document.querySelector("#search-form");i.style.display="none",r.innerHTML=`
    <div class="shared-movie-grid">
      ${e.map(t=>`
        <div class="shared-movie-card" data-movie-id="${t.id}">
          <img src="https://image.tmdb.org/t/p/w500${t.poster_path}"
               alt="${t.title}"
               class="shared-movie-image" />
          <div class="shared-movie-info">
            <h3 class="shared-movie-title">${t.title}</h3>
            <p class="shared-movie-year">${t.release_date?t.release_date.split("-")[0]:"N/A"}</p>

            ${t.rating>0?`
              <div class="shared-movie-rating">
                <span class="shared-movie-stars">${"★".repeat(t.rating)}${"☆".repeat(5-t.rating)}</span>
                <span>${t.rating}/5</span>
              </div>
            `:""}

            ${t.comment?`
              <div class="shared-movie-comment">
                "${t.comment}"
              </div>
            `:""}
          </div>
        </div>
      `).join("")}
    </div>
  `,document.getElementById("back-home").addEventListener("click",()=>{window.location.href=window.location.pathname}),document.querySelectorAll(".shared-movie-card").forEach(t=>{t.addEventListener("click",async o=>{const s=Number(o.currentTarget.dataset.movieId),a=await y(s);me(a)})})}function me(e){j(e),I(),d.innerHTML=`
    <div id="modal-content">
      <h2>${e.title} (${e.release_date?.split("-")[0]||""})</h2>
      <img src="https://image.tmdb.org/t/p/w300${e.poster_path}" />
      <p>${e.overview}</p>
      <button id="add-to-my-list">Add to My List</button>
      <button id="close-shared-modal">Close</button>
    </div>
  `,d.classList.remove("hidden"),document.querySelector("#close-shared-modal").addEventListener("click",()=>d.classList.add("hidden")),document.querySelector("#add-to-my-list").addEventListener("click",()=>{if(X().includes(e.id)){h("Movie is already in your list!","warning"),d.classList.add("hidden");return}Y(e)})}async function V(){if(f().length===0){h("Add some movies to your watched list first!","warning");return}const r=H();r?S(r):ue()}function ue(){const e=H();if(e){S(e);return}d.innerHTML=`
    <div class="dialog-content">
      <h3>What's your name?</h3>
      <p style="color: rgba(0, 0, 0, 0.7); margin-bottom: 1rem;">We'll remember this for future shares!</p>
      <input type="text" class="dialog-input" id="share-name" placeholder="Enter your name..." />
      <div class="dialog-buttons">
        <button class="dialog-btn dialog-btn-secondary" id="cancel-share">Cancel</button>
        <button class="dialog-btn dialog-btn-primary" id="confirm-share">Save & Share</button>
      </div>
    </div>
  `,d.classList.remove("hidden");const r=document.getElementById("share-name");r.focus(),document.getElementById("cancel-share").addEventListener("click",()=>{d.classList.add("hidden")}),document.getElementById("confirm-share").addEventListener("click",async()=>{const n=r.value.trim();if(!n){r.style.borderColor="#ff6b6b";return}ae(n),d.classList.add("hidden"),S(n)}),r.addEventListener("keypress",n=>{n.key==="Enter"&&document.getElementById("confirm-share").click()})}async function S(e){const r=f(),n=re(e,r),i=`${window.location.origin}${window.location.pathname}?list=${n}`;try{await navigator.clipboard.writeText(i),h("Share link copied to clipboard!","success")}catch{const t=document.createElement("textarea");t.value=i,document.body.appendChild(t),t.select(),document.execCommand("copy"),document.body.removeChild(t),h("Share link copied to clipboard!","success")}}function h(e,r="success"){const n=document.createElement("div");n.className="custom-alert",n.style.cssText=`
    position: fixed;
    top: 2rem;
    right: 2rem;
    background: ${r==="success"?"linear-gradient(45deg, #4ecdc4, #44a08d)":r==="warning"?"linear-gradient(45deg, #f093fb, #f5576c)":"linear-gradient(45deg, #ff6b6b, #ffd93d)"};
    color: white;
    padding: 1rem 1.5rem;
    border-radius: 10px;
    font-weight: 600;
    z-index: 3000;
    animation: slideInRight 0.3s ease;
    box-shadow: 0 10px 25px rgba(0, 0, 0, 0.2);
  `,n.textContent=e,document.body.appendChild(n),setTimeout(()=>{n.style.animation="slideOutRight 0.3s ease",setTimeout(()=>n.remove(),300)},3e3)}const he=new URLSearchParams(window.location.search),E=he.get("list");if(E){const e=ne(E);e&&(b=e.name,ie(e.movies))}_();I();E||Z();M.addEventListener("click",e=>{if(e.preventDefault(),f().length===0){h("You haven't watched any movies yet!","warning");return}P()});async function Z(){const e=f();l.innerHTML='<div style="text-align: center; padding: 2rem; color: rgba(0, 0, 0, 0.7);" class="loading">Finding movies you might like...</div>';try{let r=[];if(e.length>0){const i=e.slice(-3).map(async a=>{try{const c=await G(a.id),m=await Q(a.id),x=(c.results||[]).slice(0,3),L=(m.results||[]).slice(0,3);return[...x,...L]}catch{return[]}});r=(await Promise.all(i)).flat();const o=new Set(e.map(a=>a.id)),s=new Map;if(r.forEach(a=>{!o.has(a.id)&&!s.has(a.id)&&s.set(a.id,a)}),r=Array.from(s.values()).slice(0,20),r.length<15){const a=await N(),c=await T(),m=[...(a.results||[]).slice(0,10),...(c.results||[]).slice(0,10)],x=new Set([...o,...r.map($=>$.id)]),L=m.filter($=>!x.has($.id));r=[...r,...L].slice(0,20)}}else{const n=await N(),i=await T(),t=[],o=n.results||[],s=i.results||[];for(let a=0;a<20&&(a<o.length||a<s.length);a++)a<o.length&&t.push(o[a]),a<s.length&&t.length<20&&t.push(s[a]);r=t.slice(0,20)}if(r.length===0)l.innerHTML='<div style="text-align: center; padding: 2rem; color: rgba(0, 0, 0, 0.7);">Start by searching for movies to get personalized recommendations!</div>';else{const n=e.length>0?`Recommended Based on Your Last ${Math.min(e.length,3)} Movies`:"Popular & Top Rated Movies";l.innerHTML="",C(r,!1),l.innerHTML+=`
        <div style="text-align: center; margin-top: 2rem; color: rgba(0, 0, 0, 0.6); font-size: 0.9rem; font-weight: 400;">
          ${n}
        </div>
      `}}catch(r){console.error(r),l.innerHTML='<div style="text-align: center; padding: 2rem; color: rgba(255, 107, 107, 0.8);">⚠️ Unable to load recommendations. Try searching for movies!</div>'}}
