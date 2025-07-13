import useSWRSubscription from 'swr/subscription'
import './App.css'
import { MCLDash } from './components/MCLDash'


function App() {

  const { data, error } = useSWRSubscription('ws://localhost:8765/hello?sessionID=12342', (key, { next }) => {
    const socket = new WebSocket(key)
    socket.addEventListener('message', (event) => next(null, event.data))
    socket.addEventListener('error', (event) => next(event))
    return () => socket.close()
  }, { revalidateOnMount: false, revalidateIfStale: false, revalidateOnFocus: false, revalidateOnReconnect: false })
  if (error) return <div>failed to load</div>
  if (!data) return <div>loading...</div>
  let jsonData = JSON.parse(data);

  return (
    <>
      <MCLDash data={jsonData} />
    </>
  )
}

export default App
