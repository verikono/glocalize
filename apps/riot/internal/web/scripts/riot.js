// riot browser client: opens a websocket to "/ws" and reloads the page when
// the server dispatches a refresh message.
(function () {
	"use strict";

	var RECONNECT_MS = 1000;

	function connect() {
		var scheme = window.location.protocol === "https:" ? "wss:" : "ws:";
		var url = scheme + "//" + window.location.host + "/ws";
		var socket = new WebSocket(url);

		socket.addEventListener("message", function (event) {
			var message = String(event.data || "");
			if (message.indexOf("riot::refresh") === 0) {
				window.location.reload();
			}
		});

		socket.addEventListener("close", function () {
			// The server restarts on every hot rebuild, so reconnect and wait
			// for the next refresh message.
			window.setTimeout(connect, RECONNECT_MS);
		});

		socket.addEventListener("error", function () {
			socket.close();
		});
	}

	connect();
})();

const fSkillFields = `fragment SkillFields on Skill {
  id
  name
  dependencies {
    id
    name
    dependencies {
      id
      name
      dependencies {
        id
        name
        dependencies {
          id
          name
          dependencies {
            id
            name
          }
        }
      }
    }
  }
}`;

// query: graphql query, rest are cbo: { el: <element>, fn: callback function }
const LoadSkills = async (props = {}, ...cbo) => {
    // do stuff with props soon, not now.
    const { data, error: err } = await query(`{ Skills { id, name, description, url, keywords, dependencies { ...SkillFields } } } ${fSkillFields}`, props)
    if (err) {
        console.error("LoadSkills: failed getting data", err);
        for (const { el, fn } of cbo) {
            if (fn) fn(el, null, err);
        }
        return;
    }
    //onSuccessFullLoad: iterate the cbo and call them
    for (const { el, fn, } of cbo) {
        fn(el, data)
    }

}

// query posts a graphql document to /graph and returns { data, error }.
// It is synchronous in shape only: callers await the returned promise.
const query = async (document, props = {}) => {
    const query = document.replace(/[\r\n]+/g, '');
    console.log("Q", query)
    try {
        const response = await fetch("/graph", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ query, variables: props.variables || {} }),
        });
        const payload = await response.json();
        console.log("GRAPHQL RESPONSE", payload);
        if (!response.ok) {
            return { data: null, error: payload.errors || payload };
        }
        if (payload.errors && payload.errors.length) {
            return { data: payload.data, error: payload.errors };
        }
        return { data: payload.data, error: null };
    } catch (err) {
        return { data: null, error: err };
    }
}

document.addEventListener("DOMContentLoaded", evt => (async () => { await main(evt); })());

// chooser state: the catalogue the multiselect offers, the snapshot of the
// skill currently loaded into the form (so Revert has something to go back
// to), and whether the form is composing a brand new skill.
const chooserState = {
    skills: [],
    loaded: null,
    isNew: false,
};

// view state: the last catalogue payload, the hierarchy currently drawn (which
// the search narrows), the active viewstack entry, and the search string.
const viewState = {
    data: null,
    tree: null,
    view: "tree",
    search: "",
};

// viewNames is the set of viewstack entries the ?chart= search param accepts.
const viewNames = ["tree", "radial", "life", "focused-pillar"];

// viewFromParam resolves the ?chart= search param to a viewstack entry,
// falling back to the default view when it is missing or unknown.
const viewFromParam = (params) => {
    const chart = String(params.get("chart") || "").trim().toLowerCase();
    return viewNames.includes(chart) ? chart : "tree";
}

// viewRenderer resolves a viewstack entry to its renderer. It is a function
// rather than a map so the renderers, declared further down, are only looked
// up when a view is actually drawn.
const viewRenderer = (name) => {
    switch (name) {
        case "tree": return renderTreeSkillTree;
        case "life": return renderTreeOfLifeSkillTree;
        case "focused-pillar": return renderFocusedPillarSkillTree;
        default: return renderRadialSkillTree;
    }
}

const main = async (loaded,) => {
    // the graph lives in the viewstack stage, not #main itself: #main also
    // holds the resizer and the view selector, which must survive a redraw.
    const stage = document.querySelector("#charted_data");

    setLoading(stage);

    // the ?chart= search param picks the viewstack entry before anything is
    // drawn, so the first render already uses the requested view.
    const params = new URLSearchParams(window.location.search);
    viewState.view = viewFromParam(params);

    // the ?search= param restores a search, so a narrowed graph can be linked
    // to and survives a reload.
    const search = String(params.get("search") || "").trim();
    if (search) {
        viewState.search = search;
        const input = document.querySelector("#nav-search");
        if (input) input.value = search;
        setNavStatus(`search: ${search}`);
    }

    // the chooser's form is static markup, so it is wired immediately rather
    // than being replaced by a spinner and rebuilt.
    wireChooser();
    wireNavigation();
    wireResizer();
    wireViewstack();
    observeStageResize();

    await LoadSkills({}, { el: stage, fn: onLoadedSkills });

    // a ?subject=<name>&class=skill url opens straight onto that skill.
    const subject = params.get("subject");
    if (subject && params.get("class") === "skill") {
        await selectSkill(subject, { push: false });
    }
}

// renderActiveView draws the current hierarchy with the active viewstack
// renderer, then re-applies the selected-skill highlight.
const renderActiveView = () => {
    const stage = document.querySelector("#charted_data");
    if (!stage || !viewState.tree) return;

    const render = viewRenderer(viewState.view);
    render(stage, viewState.tree);

    const name = chooserState.loaded && chooserState.loaded.name;
    if (name) markActiveSkill(document.querySelector("#main"), name);
}

// observeStageResize redraws the active view when the stage changes size. It
// is installed once, so switching views does not stack observers.
const observeStageResize = () => {
    const stage = document.querySelector("#charted_data");
    if (!stage) return;

    let timer = 0;
    let last = `${stage.clientWidth}x${stage.clientHeight}`;
    const onResize = () => {
        const size = `${stage.clientWidth}x${stage.clientHeight}`;
        if (size === last) return;
        last = size;
        clearTimeout(timer);
        timer = setTimeout(() => renderActiveView(), 60);
    };

    window.addEventListener("resize", onResize);
    if (typeof ResizeObserver !== "undefined") {
        new ResizeObserver(onResize).observe(stage);
    }
}

// wireViewstack binds the graph selector in the top right of #main. A change
// redraws the graph and writes the choice back to the ?chart= search param, so
// the current view survives a reload and can be linked to.
const wireViewstack = () => {
    const select = document.querySelector("#viewstack-select");
    if (!select) return;

    select.value = viewState.view;
    select.addEventListener("change", () => {
        viewState.view = select.value;
        renderActiveView();

        const url = new URL(window.location.href);
        url.searchParams.set("chart", viewState.view);
        window.history.replaceState(window.history.state, "", url.toString());
    });
}

// wireNavigation binds the section links and the debounced skill search. A
// keystroke restarts a 1500ms timer, so the query only fires once typing has
// paused; the graph then narrows to the matches and their dependencies.
const wireNavigation = () => {
    const input = document.querySelector("#nav-search");
    if (!input) return;

    let timer = 0;
    const schedule = () => {
        clearTimeout(timer);
        timer = setTimeout(() => applySearch(input.value), 1500);
    };

    input.addEventListener("keyup", schedule);
    // the native clear button and Enter fire "search" rather than "keyup".
    input.addEventListener("search", () => {
        clearTimeout(timer);
        applySearch(input.value);
    });
}

// applySearch rebuilds the hierarchy from the catalogue, keeping only skills
// whose name matches and every dependency beneath them. The term is written
// back to the ?search= param so the narrowed graph can be linked to.
const applySearch = (value) => {
    const term = String(value || "").trim();
    viewState.search = term;

    const url = new URL(window.location.href);
    if (term) url.searchParams.set("search", term);
    else url.searchParams.delete("search");
    window.history.replaceState(window.history.state, "", url.toString());

    if (!viewState.data) return;

    viewState.tree = buildSkillHierarchy(viewState.data, term);
    setNavStatus(term ? `search: ${term}` : "");
    renderActiveView();
}

// setNavStatus writes a short message beside the search box.
const setNavStatus = (message) => {
    const status = document.querySelector("#nav-status");
    if (!status) return;
    status.textContent = message;
}

// wireResizer makes the handle at the top of #main drag the boundary between
// #chooser and #main. Only #chooser's height is set; the form's flex layout
// grows #skill-content to fill it.
const wireResizer = () => {
    const handle = document.querySelector("#chooser-resizer");
    const chooser = document.querySelector("#chooser");
    if (!handle || !chooser) return;

    let startY = 0;
    let startHeight = 0;

    const onMove = (event) => {
        const delta = event.clientY - startY;
        const max = Math.max(window.innerHeight - 160, 80);
        const height = Math.min(Math.max(startHeight + delta, 80), max);
        chooser.style.height = `${height}px`;
    };

    const onUp = () => {
        handle.classList.remove("dragging");
        window.removeEventListener("pointermove", onMove);
        window.removeEventListener("pointerup", onUp);
    };

    handle.addEventListener("pointerdown", (event) => {
        event.preventDefault();
        startY = event.clientY;
        startHeight = chooser.getBoundingClientRect().height;
        handle.classList.add("dragging");
        window.addEventListener("pointermove", onMove);
        window.addEventListener("pointerup", onUp);
    });
}

// selectSkill loads one skill into the chooser and reflects it in the url, so
// the current selection survives a reload and can be linked to.
const selectSkill = async (name, options = {}) => {
    const { push = true } = options;
    const { data, error: err } = await query(
        `query SkillByName($name: String!) { Skill(name: $name) { id name description url content keywords dependencies { id name } } }`,
        { variables: { name } },
    );
    if (err || !data || !data.Skill) {
        setStatus(err ? "failed to load skill" : "skill not found", true);
        return;
    }

    chooserState.loaded = cloneSkill(data.Skill);
    chooserState.isNew = false;
    fillForm(chooserState.loaded);

    const url = new URL(window.location.href);
    url.searchParams.set("subject", name);
    url.searchParams.set("class", "skill");
    if (push) {
        window.history.pushState({ subject: name }, "", url.toString());
    } else {
        window.history.replaceState({ subject: name }, "", url.toString());
    }

    markActiveSkill(document.querySelector("#main"), name);
}

// cloneSkill copies the fields the form edits, so the loaded snapshot is not
// mutated as the user types.
const cloneSkill = (skill) => ({
    id: skill.id || "",
    name: skill.name || "",
    description: skill.description || "",
    url: skill.url || "",
    content: skill.content || "",
    keywords: skill.keywords || "",
    dependencies: (skill.dependencies || []).map(dep => dep.id),
});

// wireChooser binds the form once: the multiselect, the new-skill button, the
// revert and save actions, and the browser back/forward buttons.
const wireChooser = () => {
    const form = document.querySelector("#skill-form");
    if (!form) return;

    const multiselect = document.querySelector("#skill-dependencies");
    const toggle = multiselect.querySelector(".multiselect-toggle");
    const menu = multiselect.querySelector(".multiselect-menu");

    toggle.addEventListener("click", (event) => {
        event.stopPropagation();
        const open = multiselect.classList.toggle("open");
        toggle.setAttribute("aria-expanded", open ? "true" : "false");
    });

    // clicking anywhere else closes the menu.
    document.addEventListener("click", (event) => {
        if (!multiselect.contains(event.target)) {
            multiselect.classList.remove("open");
            toggle.setAttribute("aria-expanded", "false");
        }
    });

    menu.addEventListener("change", () => updateMultiselectLabel());

    document.querySelector("#skill-new").addEventListener("click", () => {
        // a new skill that depends on the one currently loaded.
        const parent = chooserState.loaded;
        chooserState.isNew = true;
        fillForm({
            id: "",
            name: "",
            description: "",
            url: "",
            content: "",
            dependencies: parent && parent.id ? [parent.id] : [],
        });
        setStatus(parent ? `new skill depending on ${parent.name}` : "new skill");
        document.querySelector("#skill-name").focus();
    });

    document.querySelector("#skill-revert").addEventListener("click", () => {
        if (!chooserState.loaded) return;
        chooserState.isNew = false;
        fillForm(chooserState.loaded);
        setStatus("reverted");
    });

    form.addEventListener("submit", async (event) => {
        event.preventDefault();
        await saveSkill();
    });

    window.addEventListener("popstate", () => {
        const params = new URLSearchParams(window.location.search);

        // the viewstack entry travels in the url too, so back/forward restores
        // the view as well as the subject.
        const view = viewFromParam(params);
        if (view !== viewState.view) {
            viewState.view = view;
            const select = document.querySelector("#viewstack-select");
            if (select) select.value = view;
            renderActiveView();
        }

        const subject = params.get("subject");
        if (subject && params.get("class") === "skill") {
            selectSkill(subject, { push: false });
        }
    });
}

// fillForm writes a skill into the form and rebuilds the dependency options.
const fillForm = (skill) => {
    document.querySelector("#skill-name").value = skill.name;
    document.querySelector("#skill-description").value = skill.description;
    document.querySelector("#skill-url").value = skill.url;
    document.querySelector("#skill-content").value = skill.content;
    const keywords = document.querySelector("#skill-keywords");
    if (keywords) keywords.value = skill.keywords || "";
    renderDependencyOptions(skill.dependencies || []);
}

// renderDependencyOptions lists every skill in the catalogue as a checkbox,
// checked when it is one of the skill's dependencies. The skill being edited
// is left out: a skill cannot depend on itself.
const renderDependencyOptions = (selected) => {
    const menu = document.querySelector("#skill-dependencies .multiselect-menu");
    const chosen = new Set(selected);
    const selfId = chooserState.isNew ? "" : (chooserState.loaded && chooserState.loaded.id) || "";

    menu.textContent = "";

    const options = chooserState.skills.filter(skill => skill.id !== selfId);
    if (options.length === 0) {
        const empty = document.createElement("div");
        empty.className = "multiselect-empty";
        empty.textContent = "no skills";
        menu.appendChild(empty);
        updateMultiselectLabel();
        return;
    }

    for (const skill of options) {
        const label = document.createElement("label");
        label.className = "multiselect-option";

        const input = document.createElement("input");
        input.type = "checkbox";
        input.value = skill.id;
        input.checked = chosen.has(skill.id);

        const text = document.createElement("span");
        text.textContent = skill.name;

        label.appendChild(input);
        label.appendChild(text);
        menu.appendChild(label);
    }
    updateMultiselectLabel();
}

// updateMultiselectLabel summarises the checked dependencies on the toggle.
const updateMultiselectLabel = () => {
    const label = document.querySelector("#skill-dependencies .multiselect-label");
    const checked = selectedDependencies();
    if (checked.length === 0) {
        label.textContent = "MultiSelect <Dependencies>";
        return;
    }
    const names = checked.map(id => {
        const skill = chooserState.skills.find(s => s.id === id);
        return skill ? skill.name : id;
    });
    label.textContent = names.join(", ");
}

// selectedDependencies returns the ids of the checked dependency boxes.
const selectedDependencies = () => {
    const boxes = document.querySelectorAll("#skill-dependencies .multiselect-menu input:checked");
    return Array.from(boxes, box => box.value);
}

// saveSkill posts the form through the SaveSkill mutation and refreshes the
// tree so the new dependency edges are drawn.
const saveSkill = async () => {
    const name = document.querySelector("#skill-name").value.trim();
    const description = document.querySelector("#skill-description").value.trim();
    const url = document.querySelector("#skill-url").value.trim();
    const content = document.querySelector("#skill-content").value;
    const keywordsEl = document.querySelector("#skill-keywords");
    const keywords = keywordsEl ? keywordsEl.value.trim() : "";

    if (!name || !description || !url) {
        setStatus("name, description and url are required", true);
        return;
    }

    const input = {
        id: chooserState.isNew ? null : (chooserState.loaded && chooserState.loaded.id) || null,
        name,
        description,
        url,
        content,
        keywords,
        dependencies: selectedDependencies(),
    };

    setStatus("saving…");
    const { data, error: err } = await query(
        `mutation SaveSkill($input: SkillInput!) { SaveSkill(input: $input) { id name description url content keywords dependencies { id name } } }`,
        { variables: { input } },
    );
    if (err || !data || !data.SaveSkill) {
        setStatus("save failed", true);
        console.error("SaveSkill", err);
        return;
    }

    chooserState.loaded = cloneSkill(data.SaveSkill);
    chooserState.isNew = false;
    fillForm(chooserState.loaded);
    setStatus("saved");

    // the saved record is now the subject, so a reload reopens it.
    const savedUrl = new URL(window.location.href);
    savedUrl.searchParams.set("subject", chooserState.loaded.name);
    savedUrl.searchParams.set("class", "skill");
    window.history.replaceState({ subject: chooserState.loaded.name }, "", savedUrl.toString());

    await refreshTree();
    markActiveSkill(document.querySelector("#main"), chooserState.loaded.name);
}

// refreshTree reloads the catalogue and redraws the active view in place. It
// asks for the same recursive dependency fragment as the initial load, so a
// save does not truncate the tree to one level.
const refreshTree = async () => {
    const { data, error: err } = await query(
        `{ Skills { id, name, description, url, keywords, dependencies { ...SkillFields } } } ${fSkillFields}`,
    );
    if (err || !data) return;

    viewState.data = data;
    chooserState.skills = data.Skills || [];
    viewState.tree = buildSkillHierarchy(data, viewState.search);
    renderActiveView();
}

// markActiveSkill highlights the node for the selected skill. The node's bound
// datum is compared, not its label: labels are truncated, so a long skill name
// would never match its own shortened text.
const markActiveSkill = (el, name) => {
    if (!el) return;
    const nodes = el.querySelectorAll(".node");
    nodes.forEach(node => {
        const datum = node.__data__;
        const nodeName = datum && datum.data ? datum.data.name : "";
        node.classList.toggle("active", nodeName === name);
    });
}

// setStatus writes a short message into the form's status line.
const setStatus = (message, isError = false) => {
    const status = document.querySelector("#chooser .skill-status");
    if (!status) return;
    status.textContent = message;
    status.classList.toggle("errored", isError);
}



const setLoading = ( el ) => {
    if (!el) return;
    // empty all children of the element
    el.textContent = "";
    // append a class "load" to the element
    el.classList.add("load");
    // attach a loader spinner to the element.
    const spinner = document.createElement("div");
    spinner.className = "spinner";
    spinner.setAttribute("role", "status");
    spinner.setAttribute("aria-label", "loading");
    el.appendChild(spinner);
}

const clearLoading = (el) => {
    if (!el) return;
    el.classList.remove("load");
    const spinner = el.querySelector(".spinner");
    if (spinner) spinner.remove();
}

const onLoadedSkills = (el, data, err) => {
    console.log("onLoadedData", el, data);

    if (err || !data) {
        clearLoading(el);
        el.textContent = "";
        const p = document.createElement("p");
        p.className = "errored";
        p.textContent = "failed to load skills";
        el.appendChild(p);
        return;
    }

    // the chooser's multiselect offers the whole catalogue.
    chooserState.skills = data.Skills || [];

    viewState.data = data;
    viewState.tree = buildSkillHierarchy(data, viewState.search);
    renderActiveView();
}

// buildSkillHierarchy turns the flat Skills list (each carrying a nested
// `dependencies` subtree) into a single hierarchy rooted at a synthetic
// "skills" node. Every skill is placed at most once, so a skill that several
// others depend on appears under the first parent that claims it and cycles
// cannot recurse forever.
//
// When search is non-empty only skills whose name contains it (case
// insensitive) are kept as roots, but their dependencies are kept in full,
// recursively, so a match is never shown without the skills it relies on.
const buildSkillHierarchy = (data, search = "") => {
    const skills = (data && data.Skills) || [];
    const term = String(search || "").trim().toLowerCase();

    // index every skill we can see, including ones only reachable as a
    // dependency, so a nested stub can be replaced by its full record.
    const byId = new Map();
    const index = (skill) => {
        if (!skill || !skill.id) return;
        if (!byId.has(skill.id)) byId.set(skill.id, skill);
        for (const dep of skill.dependencies || []) index(dep);
    };
    skills.forEach(index);

    const placed = new Set();
    const build = (skill) => {
        placed.add(skill.id);
        const node = {
            id: skill.id,
            name: skill.name,
            description: skill.description || "",
            url: skill.url || "",
            keywords: skill.keywords || "",
            children: [],
        };
        for (const dep of skill.dependencies || []) {
            if (placed.has(dep.id)) continue;
            node.children.push(build(byId.get(dep.id) || dep));
        }
        return node;
    };

    const root = { id: "__root__", name: "skills", description: "", url: "", children: [] };

    // the roots are the whole catalogue, or just the search matches. A skill
    // matches on its name or on any of its keywords.
    const matches = (skill) => {
        const haystack = `${skill.name || ""} ${skill.keywords || ""}`.toLowerCase();
        return haystack.includes(term);
    };
    const roots = term
        ? skills.filter(matches)
        : skills;

    for (const skill of roots) {
        if (placed.has(skill.id)) continue;
        root.children.push(build(skill));
    }
    return root;
}

// buildDependentHierarchy turns the flat Skills list into a hierarchy that
// follows the dependency edges backwards: a skill's children are the skills
// that depend on it. The roots are the skills nothing else builds on, so the
// drawing reads left-to-right from the most fundamental skill to the most
// derived.
//
// When search is non-empty only the connected component(s) containing a match
// are kept, so a match is shown with both what it relies on and what relies on
// it.
const buildDependentHierarchy = (data, search = "") => {
    const skills = (data && data.Skills) || [];
    const term = String(search || "").trim().toLowerCase();

    // index every skill we can see, including ones only reachable as a
    // dependency, so a nested stub can be replaced by its full record.
    const byId = new Map();
    const index = (skill) => {
        if (!skill || !skill.id) return;
        if (!byId.has(skill.id)) byId.set(skill.id, skill);
        for (const dep of skill.dependencies || []) index(dep);
    };
    skills.forEach(index);

    // dependents maps a skill id to the skills that depend on it.
    const dependents = new Map();
    for (const skill of byId.values()) {
        for (const dep of skill.dependencies || []) {
            if (!dependents.has(dep.id)) dependents.set(dep.id, []);
            dependents.get(dep.id).push(skill);
        }
    }

    // visible is the set of skills to draw: everything, or the connected
    // component(s) around a match. Edges are walked in both directions so a
    // match brings in what it depends on and what depends on it.
    let visible = null;
    if (term) {
        const matches = (skill) => `${skill.name || ""} ${skill.keywords || ""}`.toLowerCase().includes(term);
        visible = new Set();
        const stack = [...byId.values()].filter(matches).map(skill => skill.id);
        while (stack.length) {
            const id = stack.pop();
            if (visible.has(id)) continue;
            visible.add(id);
            const skill = byId.get(id);
            for (const dep of (skill && skill.dependencies) || []) stack.push(dep.id);
            for (const dependent of dependents.get(id) || []) stack.push(dependent.id);
        }
    }

    const shown = (id) => !visible || visible.has(id);

    const placed = new Set();
    const build = (skill) => {
        placed.add(skill.id);
        const node = {
            id: skill.id,
            name: skill.name,
            description: skill.description || "",
            url: skill.url || "",
            keywords: skill.keywords || "",
            children: [],
        };
        for (const dependent of dependents.get(skill.id) || []) {
            if (placed.has(dependent.id) || !shown(dependent.id)) continue;
            node.children.push(build(dependent));
        }
        return node;
    };

    const root = { id: "__root__", name: "skills", description: "", url: "", keywords: "", children: [] };

    // a root is a visible skill with no visible dependency: nothing to its
    // left, so the tree starts there.
    const hasVisibleDep = (skill) => (skill.dependencies || []).some(dep => byId.has(dep.id) && shown(dep.id));
    const roots = [...byId.values()].filter(skill => shown(skill.id) && !hasVisibleDep(skill));

    for (const skill of roots) {
        if (placed.has(skill.id)) continue;
        root.children.push(build(skill));
    }

    // a cycle has no root; place whatever is left so no skill is dropped.
    for (const skill of byId.values()) {
        if (!shown(skill.id) || placed.has(skill.id)) continue;
        root.children.push(build(skill));
    }

    return root;
}

// skillTreeTruncate shortens a label so long skill names do not swamp the
// drawing. The full name is always available in the node's <title>.
const skillTreeTruncate = (name, maxChars = 22) => {
    const text = String(name || "");
    return text.length > maxChars ? text.slice(0, maxChars - 1) + "…" : text;
}

// skillTreeCaption writes the "n skills · m leaves" line into the stage.
const skillTreeCaption = (el, root) => {
    const caption = document.createElement("p");
    caption.className = "caption";
    caption.textContent = `${root.descendants().length - 1} skills · ${root.leaves().length} leaves`;
    el.appendChild(caption);
}

// skillTreeSvg clears the stage and appends a fresh svg sized to it. Every
// renderer starts this way, so switching views never leaves stale markup.
const skillTreeSvg = (el, label) => {
    clearLoading(el);
    el.textContent = "";

    return d3.select(el)
        .append("svg")
        .attr("preserveAspectRatio", "xMidYMid meet")
        .attr("role", "img")
        .attr("aria-label", label);
}

// skillTreeFitAndZoom fits the viewBox to the drawn content, then wires wheel
// zoom about the pointer and double-click to reset. The zoom is remembered on
// the element, keyed by the tree it belongs to, so a resize keeps the view.
//
// options.fit selects the initial framing: "content" (the default) fits the
// whole drawing into the stage; "width" fits the drawing's width to the stage
// and anchors it at the top, which is what a tall, row-per-item drawing wants.
const skillTreeFitAndZoom = (el, svg, treeData, options = {}) => {
    if (el.__skillTree !== treeData) el.__skillTreeZoom = null;
    el.__skillTree = treeData;

    const bounds = svg.node().getBBox();
    const pad = 8;
    const [bx, by, bw, bh] = [bounds.x - pad, bounds.y - pad, bounds.width + pad * 2, bounds.height + pad * 2];

    // the stage's aspect ratio, used to keep the viewBox from letterboxing.
    const stageWidth = Math.max(el.clientWidth || 0, 1);
    const stageHeight = Math.max(el.clientHeight || 0, 1);
    const aspect = stageHeight / stageWidth;

    const baseView = () => {
        if (options.fit !== "width") return [bx, by, bw, bh];
        // the full width, and only as much height as the stage can show at
        // that scale, anchored at the top of the drawing.
        return [bx, by, bw, bw * aspect];
    };

    // the base view's own size, which the zoom is expressed against.
    const [baseX, baseY, baseW, baseH] = baseView();
    const baseCentre = { k: 1, cx: baseX + baseW / 2, cy: baseY + baseH / 2 };

    const viewFor = (zoom) => {
        if (!zoom) return baseView();
        const { k, cx, cy } = zoom;
        return [cx - baseW / (2 * k), cy - baseH / (2 * k), baseW / k, baseH / k];
    };

    const svgNode = svg.node();
    svg.attr("viewBox", viewFor(el.__skillTreeZoom).join(" "));

    const MIN_K = 1; // never zoom out past the fitted view
    const MAX_K = 12;
    svgNode.addEventListener("wheel", (event) => {
        event.preventDefault();

        // the pointer in the svg's own coordinate system, which accounts for
        // the letterboxing preserveAspectRatio introduces.
        const ctm = svgNode.getScreenCTM();
        if (!ctm) return;
        const point = new DOMPoint(event.clientX, event.clientY).matrixTransform(ctm.inverse());

        const old = el.__skillTreeZoom || baseCentre;
        const k = Math.min(Math.max(old.k * Math.exp(-event.deltaY * 0.0015), MIN_K), MAX_K);
        if (k === old.k) return;

        const ratio = old.k / k;
        el.__skillTreeZoom = {
            k,
            cx: point.x + (old.cx - point.x) * ratio,
            cy: point.y + (old.cy - point.y) * ratio,
        };
        svg.attr("viewBox", viewFor(el.__skillTreeZoom).join(" "));
    }, { passive: false });

    svgNode.addEventListener("dblclick", (event) => {
        if (event.target !== svgNode) return;
        el.__skillTreeZoom = null;
        svg.attr("viewBox", viewFor(null).join(" "));
    });

    // middle-mouse drag pans the view, complementing the wheel zoom. The
    // pointer delta is converted into the base coordinate space through the
    // svg's own screen matrix, so panning tracks the cursor at any zoom.
    let panning = false;
    let panStart = null;

    const onPanMove = (event) => {
        if (!panning || !panStart) return;
        event.preventDefault();

        const ctm = svgNode.getScreenCTM();
        if (!ctm) return;
        const inv = ctm.inverse();

        // a screen-space delta as a base-space vector (translation excluded).
        const dx = event.clientX - panStart.x;
        const dy = event.clientY - panStart.y;
        const baseDx = inv.a * dx + inv.c * dy;
        const baseDy = inv.b * dx + inv.d * dy;

        // dragging moves the content with the cursor, so the view centre
        // travels the opposite way.
        el.__skillTreeZoom = {
            k: panStart.zoom.k,
            cx: panStart.zoom.cx - baseDx,
            cy: panStart.zoom.cy - baseDy,
        };
        svg.attr("viewBox", viewFor(el.__skillTreeZoom).join(" "));
    };

    const onPanUp = (event) => {
        if (!panning) return;
        panning = false;
        panStart = null;
        svgNode.classList.remove("panning");
        svgNode.releasePointerCapture?.(event.pointerId);
        window.removeEventListener("pointermove", onPanMove);
        window.removeEventListener("pointerup", onPanUp);
        window.removeEventListener("pointercancel", onPanUp);
    };

    svgNode.addEventListener("pointerdown", (event) => {
        // button 1 is the middle mouse button.
        if (event.button !== 1) return;
        event.preventDefault();

        panning = true;
        panStart = {
            x: event.clientX,
            y: event.clientY,
            zoom: el.__skillTreeZoom || baseCentre,
        };
        svgNode.classList.add("panning");
        svgNode.setPointerCapture?.(event.pointerId);
        window.addEventListener("pointermove", onPanMove);
        window.addEventListener("pointerup", onPanUp);
        window.addEventListener("pointercancel", onPanUp);
    });

    // suppress the browser's middle-click autoscroll and paste behaviours.
    svgNode.addEventListener("auxclick", (event) => {
        if (event.button === 1) event.preventDefault();
    });
    svgNode.addEventListener("mousedown", (event) => {
        if (event.button === 1) event.preventDefault();
    });
}

// skillTreeWireNodes binds the shared node behaviour: click loads the skill
// into the chooser, hover highlights the path from the root to the node.
const skillTreeWireNodes = (node, link) => {
    // the stylesheet makes labels transparent to the pointer so they never
    // intercept a pan; the label is still part of the node's click target, so
    // re-enable pointer events inline (which outranks the stylesheet rule).
    node.select("text").style("pointer-events", "all");

    node.append("title")
        .text(d => (d.data.description ? `${d.data.name}\n\n${d.data.description}` : d.data.name));

    node.style("cursor", d => (d.depth === 0 ? "default" : "pointer"))
        .on("click", (event, d) => {
            if (d.depth === 0) return;
            selectSkill(d.data.name);
        })
        .on("mouseenter", (event, d) => {
            const path = new Set(d.ancestors());
            link.attr("stroke-opacity", l => (path.has(l.target) ? 0.95 : 0.12))
                .attr("stroke-width", l => (path.has(l.target) ? 2.5 : 1.5));
            node.classed("active", n => path.has(n));
        })
        .on("mouseleave", () => {
            link.attr("stroke-opacity", 0.35).attr("stroke-width", 1.5);
            node.classed("active", false);

            // the hover highlight is transient; the selected skill's highlight
            // is not, so restore it once the pointer leaves.
            const name = chooserState.loaded && chooserState.loaded.name;
            if (name) markActiveSkill(document.querySelector("#main"), name);
        });
}

// renderRadialSkillTree draws the hierarchy as a radial tree.
// https://observablehq.com/@d3/radial-tree/2
const renderRadialSkillTree = (el, treeData) => {
    if (!el || !treeData) return;

    if (typeof d3 === "undefined") {
        renderD3Missing(el);
        return;
    }

    const width = Math.max(el.clientWidth || 0, 320);
    const height = Math.max(el.clientHeight || 0, 320);
    const fontSize = 9;

    // The layout is computed in a square of this radius; the viewBox is fitted
    // to the drawn content afterwards, so the tree fills whatever space the
    // container has regardless of how long the labels turn out to be.
    const radius = Math.max(Math.min(width, height) / 2 - 40, 60);

    const root = d3.hierarchy(treeData);
    d3.tree()
        .size([2 * Math.PI, radius])
        .separation((a, b) => (a.parent === b.parent ? 1 : 2) / a.depth)(root);

    // The skill graph is mostly flat (many roots, few levels), so a linear
    // radius per depth would bunch every root near the centre. Push the inner
    // rings outward instead, leaving the outer ring at the full radius.
    const maxDepth = d3.max(root.descendants(), d => d.depth) || 1;
    root.each(d => {
        d.y = d.depth === 0 ? 0 : radius * Math.pow(d.depth / maxDepth, 0.55);
    });

    const svg = skillTreeSvg(el, "radial tree of skills and their dependencies");

    // links: a radial curve from each parent to each child.
    const link = svg.append("g")
        .attr("class", "links")
        .attr("fill", "none")
        .attr("stroke", "currentColor")
        .attr("stroke-opacity", 0.35)
        .attr("stroke-width", 1.5)
        .selectAll("path")
        .data(root.links())
        .join("path")
        .attr("d", d3.linkRadial().angle(d => d.x).radius(d => d.y));

    // nodes: each group is rotated to its angle and pushed out to its radius.
    const node = svg.append("g")
        .attr("class", "nodes")
        .selectAll("g")
        .data(root.descendants())
        .join("g")
        .attr("class", "node")
        .attr("transform", d => `rotate(${d.x * 180 / Math.PI - 90}) translate(${d.y},0)`);

    // an invisible, generous hit area so a node is easy to hover and click
    // even though the visible circle is only a few pixels across.
    node.append("circle")
        .attr("class", "hit")
        .attr("r", 10)
        .attr("fill", "transparent");

    node.append("circle")
        .attr("r", d => (d.depth === 0 ? 5 : d.children ? 4 : 3))
        .attr("fill", d => (d.children ? "currentColor" : "#8a8a8a"));

    node.append("text")
        .attr("dy", "0.31em")
        .attr("x", d => (d.x < Math.PI === !d.children ? 6 : -6))
        .attr("text-anchor", d => (d.x < Math.PI === !d.children ? "start" : "end"))
        .attr("transform", d => (d.x >= Math.PI ? "rotate(180)" : null))
        .attr("paint-order", "stroke")
        .attr("stroke", "#fff")
        .attr("stroke-width", 3)
        .attr("stroke-linejoin", "round")
        .attr("font-size", fontSize)
        .attr("font-family", "sans-serif")
        .text(d => (d.depth === 0 ? "" : skillTreeTruncate(d.data.name)));

    // the synthetic root is unlabelled, so name the chart instead.
    skillTreeCaption(el, root);

    skillTreeWireNodes(node, link);
    skillTreeFitAndZoom(el, svg, treeData);
}

// renderTreeSkillTree draws the hierarchy as an indented left-to-right tree,
// one row per node, with step-before links. The hierarchy follows the
// dependency edges backwards, so a skill's dependencies sit to its left and
// the skills that depend on it sit to its right.
// https://observablehq.com/@observablehq/plot-custom-tree-layout
const renderTreeSkillTree = (el) => {
    if (!el) return;

    if (typeof d3 === "undefined") {
        renderD3Missing(el);
        return;
    }

    const hierarchy = buildDependentHierarchy(viewState.data, viewState.search);

    const fontSize = 10;
    const rowHeight = 20;
    const columnWidth = 190;

    const root = d3.hierarchy(hierarchy);

    // the tidy layout: depth becomes the column and siblings stack vertically,
    // so a chain of skills reads as one straight horizontal row rather than a
    // staircase. nodeSize is [row, column] because the tree grows sideways.
    d3.tree().nodeSize([rowHeight, columnWidth])(root);

    const svg = skillTreeSvg(el, "indented tree of skills and their dependencies");

    // the synthetic root is not drawn: it would sit a whole column to the left
    // of the first real skill and add a long empty link. Columns are shifted
    // back by one so the first real skill starts at the left edge.
    const posX = (d) => (d.depth - 1) * columnWidth;

    // links: a step-before path, horizontal to the child's column then
    // vertical to its row.
    const link = svg.append("g")
        .attr("class", "links")
        .attr("fill", "none")
        .attr("stroke", "currentColor")
        .attr("stroke-opacity", 0.35)
        .attr("stroke-width", 1.5)
        .selectAll("path")
        .data(root.links().filter(l => l.source.depth > 0))
        .join("path")
        .attr("d", d => `M${posX(d.source)},${d.source.x}H${posX(d.target)}V${d.target.x}`);

    const node = svg.append("g")
        .attr("class", "nodes")
        .selectAll("g")
        .data(root.descendants().filter(d => d.depth > 0))
        .join("g")
        .attr("class", "node")
        .attr("transform", d => `translate(${posX(d)},${d.x})`);

    node.append("circle")
        .attr("class", "hit")
        .attr("r", 9)
        .attr("fill", "transparent");

    node.append("circle")
        .attr("r", d => (d.children ? 3.5 : 2.5))
        .attr("fill", d => (d.children ? "currentColor" : "#8a8a8a"));

    node.append("text")
        .attr("dy", "0.31em")
        .attr("x", 8)
        .attr("paint-order", "stroke")
        .attr("stroke", "#fff")
        .attr("stroke-width", 3)
        .attr("stroke-linejoin", "round")
        .attr("font-size", fontSize)
        .attr("font-family", "sans-serif")
        .text(d => skillTreeTruncate(d.data.name, 30));

    skillTreeCaption(el, root);

    skillTreeWireNodes(node, link);
    // the zoom is keyed by the catalogue hierarchy, not the dependent one
    // rebuilt here, so a resize redraw keeps the user's pan and zoom.
    skillTreeFitAndZoom(el, svg, viewState.tree);
}

// renderTreeOfLifeSkillTree draws the hierarchy as a radial cluster whose
// rings are spaced by cumulative branch length, coloured by top-level branch.
// https://observablehq.com/@d3/tree-of-life
const renderTreeOfLifeSkillTree = (el, treeData) => {
    if (!el || !treeData) return;

    if (typeof d3 === "undefined") {
        renderD3Missing(el);
        return;
    }

    const width = Math.max(el.clientWidth || 0, 320);
    const height = Math.max(el.clientHeight || 0, 320);
    const fontSize = 9;

    const outerRadius = Math.max(Math.min(width, height) / 2 - 30, 80);
    const innerRadius = Math.max(outerRadius - 90, 40);

    const root = d3.hierarchy(treeData);

    // every skill counts as one unit of branch length, so a deeper chain
    // reaches further out than a shallow one. The synthetic root contributes
    // no length of its own, so it sits at the centre.
    const maxLength = (d) => (d.depth === 0 ? 0 : 1) + (d.children ? d3.max(d.children, maxLength) : 0);
    const total = maxLength(root) || 1;

    d3.cluster()
        .size([360, innerRadius])
        .separation(() => 1)(root);

    // setRadius walks the tree summing the branch lengths, scaling the total
    // so the deepest node lands on the outer radius.
    const setRadius = (d, y0) => {
        d.radius = (y0 += d.depth === 0 ? 0 : 1) * (outerRadius / total);
        if (d.children) d.children.forEach(child => setRadius(child, y0));
    };
    setRadius(root, 0);

    // colour is inherited from the top-level branch, so each root skill's
    // whole subtree shares a hue.
    const color = d3.scaleOrdinal(d3.schemeCategory10);
    const setColor = (d) => {
        d.color = d.depth === 1 ? color(d.data.name) : d.parent ? d.parent.color : null;
        if (d.children) d.children.forEach(setColor);
    };
    setColor(root);

    const svg = skillTreeSvg(el, "tree of life of skills and their dependencies");

    // a step link: an arc along the parent's radius, then a radial line out.
    const linkStep = (startAngle, startRadius, endAngle, endRadius) => {
        const a0 = (startAngle - 90) / 180 * Math.PI;
        const a1 = (endAngle - 90) / 180 * Math.PI;
        const c0 = Math.cos(a0);
        const s0 = Math.sin(a0);
        const c1 = Math.cos(a1);
        const s1 = Math.sin(a1);
        return "M" + startRadius * c0 + "," + startRadius * s0
            + (endAngle === startAngle ? "" : "A" + startRadius + "," + startRadius + " 0 0 " + (endAngle > startAngle ? 1 : 0) + " " + startRadius * c1 + "," + startRadius * s1)
            + "L" + endRadius * c1 + "," + endRadius * s1;
    };

    const link = svg.append("g")
        .attr("class", "links")
        .attr("fill", "none")
        .attr("stroke", "currentColor")
        .attr("stroke-opacity", 0.35)
        .attr("stroke-width", 1.5)
        .selectAll("path")
        .data(root.links())
        .join("path")
        .attr("d", d => linkStep(d.source.x, d.source.radius, d.target.x, d.target.radius));

    const node = svg.append("g")
        .attr("class", "nodes")
        .selectAll("g")
        .data(root.descendants())
        .join("g")
        .attr("class", "node")
        .attr("transform", d => `rotate(${d.x - 90}) translate(${d.radius},0)`);

    node.append("circle")
        .attr("class", "hit")
        .attr("r", 9)
        .attr("fill", "transparent");

    node.append("circle")
        .attr("r", d => (d.depth === 0 ? 4 : d.children ? 3.5 : 2.5))
        .attr("fill", d => d.color || "currentColor");

    // labels sit outside the outermost ring, rotated to read outward.
    node.append("text")
        .attr("dy", "0.31em")
        .attr("x", d => (d.x < 180 ? 6 : -6))
        .attr("text-anchor", d => (d.x < 180 ? "start" : "end"))
        .attr("transform", d => (d.x >= 180 ? "rotate(180)" : null))
        .attr("paint-order", "stroke")
        .attr("stroke", "#fff")
        .attr("stroke-width", 3)
        .attr("stroke-linejoin", "round")
        .attr("font-size", fontSize)
        .attr("font-family", "sans-serif")
        .text(d => (d.depth === 0 ? "" : skillTreeTruncate(d.data.name)));

    skillTreeCaption(el, root);

    skillTreeWireNodes(node, link);
    skillTreeFitAndZoom(el, svg, treeData);
}

// buildDependentIndex indexes the catalogue by id and, separately, records for
// every skill the skills that depend on it. The dependency edges point from a
// skill to what it needs; the focused-pillar view walks them backwards.
const buildDependentIndex = (skills) => {
    const byId = new Map();
    const index = (skill) => {
        if (!skill || !skill.id) return;
        if (!byId.has(skill.id)) byId.set(skill.id, skill);
        for (const dep of skill.dependencies || []) index(dep);
    };
    skills.forEach(index);

    const dependents = new Map();
    for (const skill of byId.values()) {
        for (const dep of skill.dependencies || []) {
            if (!dependents.has(dep.id)) dependents.set(dep.id, []);
            dependents.get(dep.id).push(skill);
        }
    }
    return { byId, dependents };
}

// renderFocusedPillarSkillTree draws one row per skill: the skill's name in a
// central column, and to its right a tree of everything that depends on it,
// each level one column further right. Rows are stacked vertically and the
// stage scrolls, so every skill gets a full row of its own.
const renderFocusedPillarSkillTree = (el, treeData) => {
    if (!el || !treeData) return;

    if (typeof d3 === "undefined") {
        renderD3Missing(el);
        return;
    }

    const skills = (viewState.data && viewState.data.Skills) || [];
    const term = String(viewState.search || "").trim().toLowerCase();
    const { dependents } = buildDependentIndex(skills);

    const fontSize = 10;
    const nodeRowHeight = 22;
    const columnWidth = 170;
    const nameColumnWidth = 260;
    const rowPadding = 16;
    const minRowHeight = 36;
    // the caption is absolutely positioned over the stage, so the rows start
    // below it rather than underneath it.
    const headerHeight = 26;

    const matches = (skill) => {
        const haystack = `${skill.name || ""} ${skill.keywords || ""}`.toLowerCase();
        return haystack.includes(term);
    };
    const rows = term ? skills.filter(matches) : skills;

    // buildRow turns one skill into a dependents hierarchy. The skill itself is
    // the synthetic root (it is already drawn in the name column), so only its
    // dependents and their dependents become nodes. `seen` keeps a cycle from
    // recursing forever.
    const buildRow = (skill) => {
        const seen = new Set([skill.id]);
        const walk = (node) => {
            const children = [];
            for (const dependent of dependents.get(node.id) || []) {
                if (seen.has(dependent.id)) continue;
                seen.add(dependent.id);
                children.push(walk(dependent));
            }
            return {
                id: node.id,
                name: node.name,
                description: node.description || "",
                url: node.url || "",
                keywords: node.keywords || "",
                children,
            };
        };
        return { id: "__row__", name: "", description: "", url: "", keywords: "", children: walk(skill).children };
    };

    // lay every row out first: the tree's vertical extent decides the row's
    // height, and the rows are then stacked.
    let cursorY = headerHeight;
    let maxDepth = 0;
    const layouts = rows.map((skill) => {
        const root = d3.hierarchy(buildRow(skill));
        d3.tree().nodeSize([nodeRowHeight, columnWidth])(root);

        const xs = root.descendants().map(d => d.x);
        const minX = Math.min(...xs);
        const maxX = Math.max(...xs);
        const height = Math.max(maxX - minX + nodeRowHeight + rowPadding, minRowHeight);

        maxDepth = Math.max(maxDepth, d3.max(root.descendants(), d => d.depth) || 0);

        const layout = { skill, root, minX, maxX, height, top: cursorY };
        cursorY += height;
        return layout;
    });

    const contentWidth = nameColumnWidth + maxDepth * columnWidth + 200;
    const contentHeight = Math.max(cursorY, 1);

    const svg = skillTreeSvg(el, "focused pillar of skills and their dependents");

    // the drawing is laid out in its own pixel space and then fitted into the
    // stage by the shared zoom/pan helper, exactly like the other views. #main
    // stays overflow:hidden, so the wheel and middle-drag gestures are the only
    // way to move around a drawing taller than the stage.
    svg.attr("viewBox", `0 0 ${contentWidth} ${contentHeight}`);

    // the rule that separates the name column from the dependents trees.
    svg.append("line")
        .attr("class", "pillar-rule")
        .attr("x1", nameColumnWidth)
        .attr("x2", nameColumnWidth)
        .attr("y1", 0)
        .attr("y2", contentHeight)
        .attr("stroke", "currentColor")
        .attr("stroke-opacity", 0.18)
        .attr("stroke-width", 1);

    const allNodes = [];
    const allLinks = [];

    for (const layout of layouts) {
        const { skill, root, minX, maxX, height, top } = layout;
        const row = svg.append("g")
            .attr("class", "pillar-row")
            .attr("transform", `translate(0,${top})`);

        // the name column: the skill itself, centred in its own column.
        const name = row.append("text")
            .attr("class", "pillar-name")
            .attr("x", nameColumnWidth / 2)
            .attr("y", height / 2)
            .attr("dy", "0.31em")
            .attr("text-anchor", "middle")
            .attr("font-size", fontSize + 1)
            .attr("font-family", "sans-serif")
            .style("cursor", "pointer")
            .text(skillTreeTruncate(skill.name, 34));

        name.append("title").text(skill.description ? `${skill.name}\n\n${skill.description}` : skill.name);
        name.on("click", () => selectSkill(skill.name));

        // the tree is centred vertically in the row, and its first level is
        // pulled back to the name column's edge.
        const yOffset = height / 2 - (minX + maxX) / 2;
        const posX = (d) => nameColumnWidth + (d.y - columnWidth);
        const posY = (d) => yOffset + d.x;

        const links = root.links().filter(l => l.source.depth > 0);
        const inner = row.append("g")
            .attr("class", "links")
            .attr("fill", "none")
            .attr("stroke", "currentColor")
            .attr("stroke-opacity", 0.35)
            .attr("stroke-width", 1.5)
            .selectAll("path")
            .data(links)
            .join("path")
            .attr("d", d => `M${posX(d.source)},${posY(d.source)}H${posX(d.target)}V${posY(d.target)}`);

        // the first level is joined straight to the name column, since the
        // skill itself is not drawn as a node. The data is shaped like a link
        // so the shared hover highlight treats it the same as the rest.
        const first = row.append("g")
            .attr("class", "links")
            .attr("fill", "none")
            .attr("stroke", "currentColor")
            .attr("stroke-opacity", 0.35)
            .attr("stroke-width", 1.5)
            .selectAll("path")
            .data((root.children || []).map(child => ({ source: root, target: child })))
            .join("path")
            .attr("d", d => `M${nameColumnWidth},${height / 2}H${posX(d.target)}V${posY(d.target)}`);

        // one selection so the hover highlight covers both link groups.
        const link = inner.merge(first);

        const node = row.append("g")
            .attr("class", "nodes")
            .selectAll("g")
            .data(root.descendants().filter(d => d.depth > 0))
            .join("g")
            .attr("class", "node")
            .attr("transform", d => `translate(${posX(d)},${posY(d)})`);

        node.append("circle")
            .attr("class", "hit")
            .attr("r", 9)
            .attr("fill", "transparent");

        node.append("circle")
            .attr("r", d => (d.children ? 3.5 : 2.5))
            .attr("fill", d => (d.children ? "currentColor" : "#8a8a8a"));

        node.append("text")
            .attr("dy", "0.31em")
            .attr("x", 8)
            .attr("paint-order", "stroke")
            .attr("stroke", "#fff")
            .attr("stroke-width", 3)
            .attr("stroke-linejoin", "round")
            .attr("font-size", fontSize)
            .attr("font-family", "sans-serif")
            .text(d => skillTreeTruncate(d.data.name, 26));

        skillTreeWireNodes(node, link);

        allNodes.push(node);
        allLinks.push(link);
    }

    // the caption counts the rows and the dependent edges drawn beneath them.
    const caption = document.createElement("p");
    caption.className = "caption";
    const edges = allLinks.reduce((total, link) => total + link.size(), 0);
    caption.textContent = `${rows.length} skills · ${edges} dependent links`;
    el.appendChild(caption);

    // the zoom is keyed by the catalogue hierarchy, not the dependent one
    // rebuilt here, so a resize redraw keeps the user's pan and zoom. The
    // drawing is tall and narrow, so it is fitted to the stage's width and
    // anchored at the top rather than shrunk to fit the whole height.
    skillTreeFitAndZoom(el, svg, viewState.tree, { fit: "width" });
}

// renderD3Missing reports a failed d3 load in place of a graph.
const renderD3Missing = (el) => {
    clearLoading(el);
    el.textContent = "";
    const p = document.createElement("p");
    p.className = "errored";
    p.textContent = "d3 failed to load";
    el.appendChild(p);
}