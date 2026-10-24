// The page module of Greet, on the plugin SDK: its views in Solid with the UI kit, built by the
// preset into one ES module (web/index.js). The IDE runs it in a sandboxed frame (no network,
// no storage, no access to the page) and gives it `ide`, the sandbox API.
import { createSignal } from 'solid-js'
import { useIde, view } from '@antennae-space/plugin-sdk/solid'
import { Button, Message, View } from '@antennae-space/plugin-sdk/ui'
import './greet.css'

// The tool of the rail: a greeting, the count of the files by the backend, a note to open.
function Greet() {
  const ide = useIde()
  const [answer, setAnswer] = createSignal<{ text: string; error?: boolean }>()
  const count = async () => {
    try {
      const n = await ide.rpc<number>('greet.count')
      setAnswer({ text: ide.t('{n} entries at the root of the project', { n }) })
    } catch (e) {
      setAnswer({ text: e instanceof Error ? e.message : String(e), error: true })
    }
  }
  return (
    <View>
      <div class="greet">
        <strong id="hello">{ide.t('Hello from Greet')}</strong>
        <Button id="count" icon="search" onClick={count}>{ide.t('Count the files of the project')}</Button>
        {answer() && <Message kind={answer()!.error ? 'error' : 'info'}><span id="answer">{answer()!.text}</span></Message>}
        <Button id="note" icon="file" onClick={() => ide.openView('note', { opened: Date.now() }, 'Greet')}>{ide.t('Open a note')}</Button>
      </div>
    </View>
  )
}

// A tab view: its state is kept with the tab, across reloads.
function Note() {
  const ide = useIde()
  const opened = (ide.view?.state as { opened?: number } | undefined)?.opened
  return (
    <div class="greet">
      <p>{ide.t('A note of Greet, kept with the tab.')}</p>
      <p class="ide-muted">{opened ? new Date(opened).toLocaleString(ide.lang ?? undefined) : ''}</p>
    </div>
  )
}

ide.register({
  views: { greet: view(Greet), note: view(Note) },
  commands: {
    hello: (ide) => ide.toast(ide.t('Hello!')),
  },
})
