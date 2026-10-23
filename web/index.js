// The page module of Greet: one plain ES module, no build. The IDE runs it in a sandboxed
// frame (no network, no storage, no access to the page) and gives it `ide`, the sandbox API.
// It registers its views (a tool, a tab view) and its commands.

const style = `
  .greet { padding: 12px; display: flex; flex-direction: column; gap: 8px; }
  .greet button { align-self: flex-start; background: var(--bg-3, #333); color: var(--fg, #eee);
    border: 1px solid var(--line, #555); border-radius: 4px; padding: 4px 10px; cursor: pointer; }
  .greet .muted { color: var(--fg-muted, #999); }
`

function frame(el, html) {
  el.innerHTML = `<style>${style}</style><div class="greet">${html}</div>`
  return el.querySelector('.greet')
}

ide.register({
  views: {
    // The tool of the rail: a greeting, the count of the files by the backend, a note to open.
    greet: (el, ide) => {
      const root = frame(
        el,
        `<strong id="hello">${ide.t('Hello from Greet')}</strong>
         <button id="count">${ide.t('Count the files of the project')}</button>
         <span id="answer" class="muted"></span>
         <button id="note">${ide.t('Open a note')}</button>`,
      )
      root.querySelector('#count').onclick = async () => {
        const out = root.querySelector('#answer')
        try {
          const n = await ide.rpc('greet.count')
          out.textContent = ide.t('{n} entries at the root of the project', { n })
        } catch (e) {
          out.textContent = e.message
        }
      }
      root.querySelector('#note').onclick = () => ide.openView('note', { opened: Date.now() }, 'Greet')
    },
    // A tab view: its state is kept with the tab, across reloads.
    note: (el, ide) => {
      const opened = ide.view.state && ide.view.state.opened
      frame(el, `<p>${ide.t('A note of Greet, kept with the tab.')}</p><p class="muted">${opened ? new Date(opened).toLocaleString(ide.lang) : ''}</p>`)
    },
  },
  commands: {
    hello: (ide) => ide.toast(ide.t('Hello!')),
  },
})
