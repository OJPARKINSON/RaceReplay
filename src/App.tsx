import useSWRSubscription from 'swr/subscription'
import './App.css'
import { MCLDash } from './components/MCLDash'

function App() {

  const { data, error } = useSWRSubscription('ws://localhost:8765/hello?sessionID=12342', (key, { next }) => {
    const socket = new WebSocket(key)
    socket.addEventListener('message', (event) => next(null, event.data))
    socket.addEventListener('error', (event) => next(event))
    return () => socket.close()
  })
  if (error) return <div>failed to load</div>
  if (!data) return <div>loading...</div>
  return <div>hello {data}!</div>
  return (
    <>
      <MCLDash />
    </>
  )
}

export default App
