import { DashboardPage } from "./pages/DashboardPage"
import { Toaster } from "sonner"

function App() {
  return (
    <>
      <DashboardPage />
      <Toaster position="bottom-right" />
    </>
  )
}

export default App