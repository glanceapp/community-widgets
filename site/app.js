const REPO_URL = "https://github.com/glanceapp/community-widgets";
const WIDGETS_PATH = "../widgets";
const NEW_WIDGETS_COUNT = 5;

const TYPES = [
  { id: "all", label: "All" },
  { id: "custom-api", label: "Custom API" },
  { id: "extension", label: "Extensions" },
];

const $ = (id) => document.getElementById(id);

const state = { query: "", type: "all", category: "All", noKey: false, sort: "title" };
let widgets = [];
let detailRequest = 0;

function el(tag, props = {}, ...children) {
  const node = Object.assign(document.createElement(tag), props);
  node.append(...children.filter((child) => child != null && child !== false));
  return node;
}

function capitalize(text) {
  return text.charAt(0).toUpperCase() + text.slice(1);
}

function formatDate(value) {
  return new Date(value).toLocaleDateString(undefined, { year: "numeric", month: "short", day: "numeric" });
}

async function init() {
  const response = await fetch("data.json");
  widgets = await response.json();

  const newest = widgets
    .filter((w) => w.time_added)
    .sort((a, b) => b.time_added.localeCompare(a.time_added))
    .slice(0, NEW_WIDGETS_COUNT);

  for (const w of widgets) {
    w.description = capitalize(w.description);
    w.isNew = newest.includes(w);
    w.searchText = [w.title, w.description, w.author, w.directory, w.category, ...(w.env_vars || [])]
      .join(" ")
      .toLowerCase();
  }

  $("search").addEventListener("input", (e) => update({ query: e.target.value }));
  $("sort").addEventListener("change", (e) => update({ sort: e.target.value }));
  $("no-key").addEventListener("change", (e) => update({ noKey: e.target.checked }));
  $("reset").addEventListener("click", () => update({ query: "", type: "all", category: "All", noKey: false }));

  document.addEventListener("keydown", (e) => {
    if (e.key === "/" && !e.target.closest("input, textarea, select") && !$("detail").open) {
      e.preventDefault();
      $("search").focus();
    }
  });

  const detail = $("detail");
  $("detail-close").addEventListener("click", closeDetail);
  detail.addEventListener("click", (e) => e.target === detail && closeDetail());
  // Covers closing with the Escape key
  detail.addEventListener("close", clearDetailHash);

  window.addEventListener("hashchange", route);

  render();
  route();
}

function update(changes) {
  Object.assign(state, changes);
  $("search").value = state.query;
  $("no-key").checked = state.noKey;
  render();
}

function matches(w, { ignoreCategory = false } = {}) {
  if (state.type !== "all" && w.type !== state.type) return false;
  if (state.noKey && (w.type !== "custom-api" || w.needs_key)) return false;
  if (!ignoreCategory && state.category !== "All" && w.category !== state.category) return false;

  const terms = state.query.toLowerCase().split(/\s+/).filter(Boolean);
  return terms.every((term) => w.searchText.includes(term));
}

function compare(a, b) {
  if (state.sort !== "title") {
    const field = state.sort === "added" ? "time_added" : "time_updated";
    const byDate = (b[field] || "").localeCompare(a[field] || "");
    if (byDate !== 0) return byDate;
  }

  return a.title.localeCompare(b.title, undefined, { sensitivity: "base" });
}

function render() {
  renderTypes();
  renderCategories();

  const visible = widgets.filter((w) => matches(w)).sort(compare);

  $("count").textContent = `Showing ${visible.length} of ${widgets.length} widgets`;
  $("empty").hidden = visible.length > 0;
  $("grid").replaceChildren(...visible.map(renderCard));
}

function filterButton(label, count, pressed, onClick) {
  const button = el("button", { type: "button", onclick: onClick }, label, el("span", {}, String(count)));
  button.setAttribute("aria-pressed", String(pressed));
  return button;
}

function renderTypes() {
  $("types").replaceChildren(
    ...TYPES.map(({ id, label }) => {
      const count = widgets.filter((w) => id === "all" || w.type === id).length;
      return filterButton(label, count, state.type === id, () => update({ type: id }));
    })
  );
}

function renderCategories() {
  const candidates = widgets.filter((w) => matches(w, { ignoreCategory: true }));
  const counts = new Map();
  for (const w of candidates) {
    counts.set(w.category, (counts.get(w.category) || 0) + 1);
  }

  const names = [...new Set(widgets.map((w) => w.category))]
    .filter((name) => counts.has(name) || name === state.category)
    .sort((a, b) => (a === "Other") - (b === "Other") || a.localeCompare(b));

  $("categories").replaceChildren(
    filterButton("All", candidates.length, state.category === "All", () => update({ category: "All" })),
    ...names.map((name) =>
      filterButton(name, counts.get(name) || 0, state.category === name, () => update({ category: name }))
    )
  );
}

function renderCard(w) {
  const isExtension = w.type === "extension";

  const stage = w.preview
    ? el(
        "div",
        { className: "card-stage" },
        el("img", { src: `${WIDGETS_PATH}/${w.directory}/${w.preview}`, alt: "", loading: "lazy" })
      )
    : el("div", { className: "card-stage card-stage-empty" }, "Preview available on the project page");

  const variables = w.env_vars?.length || 0;

  const card = el(
    "a",
    {
      className: "card",
      href: isExtension ? w.url : `#/widgets/${w.directory}`,
    },
    stage,
    el(
      "div",
      { className: "card-body" },
      el("h3", { className: "card-title" }, w.title),
      w.author && el("span", { className: "card-author" }, `by @${w.author}`),
      el("p", { className: "card-description" }, w.description),
      el(
        "div",
        { className: "tags" },
        w.isNew && el("span", { className: "tag tag-new" }, "New"),
        isExtension && el("span", { className: "tag tag-accent" }, "Extension ↗"),
        el("span", { className: "tag" }, w.category),
        w.needs_key && el("span", { className: "tag" }, "API key"),
        variables > 0 && el("span", { className: "tag" }, `${variables} variable${variables === 1 ? "" : "s"}`)
      )
    )
  );

  if (isExtension) {
    card.target = "_blank";
    card.rel = "noopener";
  }

  return card;
}

function route() {
  const match = location.hash.match(/^#\/widgets\/(.+)$/);
  const widget = match && widgets.find((w) => w.directory === decodeURIComponent(match[1]));

  if (widget) {
    openDetail(widget);
  } else if ($("detail").open) {
    $("detail").close();
  }
}

function clearDetailHash() {
  if (location.hash) history.replaceState(null, "", location.pathname + location.search);
}

function closeDetail() {
  $("detail").close();
  clearDetailHash();
}

async function openDetail(w) {
  const request = ++detailRequest;
  const sourceUrl = `${REPO_URL}/tree/main/widgets/${w.directory}`;
  const body = $("detail-body");

  $("detail-title").textContent = w.title;
  $("detail-meta").replaceChildren(
    w.author ? el("a", { href: `https://github.com/${w.author}`, target: "_blank", rel: "noopener" }, `@${w.author}`) : "",
    ` · ${w.category} · added ${formatDate(w.time_added)} · updated ${formatDate(w.time_updated)} · `,
    el("a", { href: sourceUrl, target: "_blank", rel: "noopener" }, "View on GitHub")
  );

  $("detail-setup").replaceChildren(
    ...(w.env_vars?.length ? ["Variables to set: ", ...w.env_vars.map((name) => el("code", {}, name))] : [])
  );

  body.replaceChildren(el("p", {}, "Loading…"));
  if (!$("detail").open) $("detail").showModal();
  $("detail").scrollTop = 0;

  let markdown;
  try {
    const response = await fetch(`${WIDGETS_PATH}/${w.directory}/README.md`);
    if (!response.ok) throw new Error(response.statusText);
    markdown = await response.text();
  } catch {
    if (request !== detailRequest) return;
    body.replaceChildren(
      el("p", {}, "Couldn't load this widget's README. ", el("a", { href: sourceUrl }, "Open it on GitHub"), ".")
    );
    return;
  }

  if (request !== detailRequest) return;
  body.replaceChildren(renderReadme(markdown, w));
}

function isRelativeUrl(url) {
  return Boolean(url) && !/^([a-z][a-z0-9+.-]*:|\/\/|#)/i.test(url);
}

function renderReadme(markdown, w) {
  // Fall back to the raw markdown when the CDN scripts are unavailable
  if (!window.marked || !window.DOMPurify) {
    return el("div", { className: "code-block" }, el("pre", {}, markdown));
  }

  const fragment = DOMPurify.sanitize(marked.parse(markdown), { RETURN_DOM_FRAGMENT: true });

  // The dialog header already shows the title
  if (fragment.firstElementChild?.tagName === "H1") {
    fragment.firstElementChild.remove();
  }

  for (const img of fragment.querySelectorAll("img")) {
    const src = img.getAttribute("src");
    if (isRelativeUrl(src)) img.src = `${WIDGETS_PATH}/${w.directory}/${src}`;
    img.loading = "lazy";
  }

  for (const link of fragment.querySelectorAll("a[href]")) {
    const href = link.getAttribute("href");
    if (href.startsWith("#")) continue;
    if (isRelativeUrl(href)) link.href = `${REPO_URL}/blob/main/widgets/${w.directory}/${href}`;
    link.target = "_blank";
    link.rel = "noopener";
  }

  for (const pre of fragment.querySelectorAll("pre")) {
    const button = el("button", { type: "button" }, "Copy");
    button.addEventListener("click", async () => {
      try {
        await navigator.clipboard.writeText(pre.textContent);
        button.textContent = "Copied";
      } catch {
        button.textContent = "Copy failed";
      }
      setTimeout(() => (button.textContent = "Copy"), 1500);
    });

    const wrapper = el("div", { className: "code-block" });
    pre.replaceWith(wrapper);
    wrapper.append(pre, button);
  }

  return fragment;
}

init();
