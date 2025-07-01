import { BrowserRouter, Routes, Route } from 'react-router-dom';
import './App.css'
import Dashboard from './components/Dashboard.jsx'
import Signin from './components/Signin.jsx'
import Login from './components/login.jsx';
import Admin from './components/Admin.jsx'
import Users from './components/Users.jsx'
import Summary from './components/Summary.jsx'
import SummeryUser from './components/SummeryUser.jsx'
function App() {
  

  return (
    <> <BrowserRouter>
      <Routes>
        <Route path="/" element={<Login />} />
        <Route path="/Signin" element={<Signin />} />
        <Route path="/Dashboard" element={<Dashboard />} />
        <Route path="/Admin" element={<Admin />} />
        <Route path="/Users" element={<Users />} />
        <Route path="/Summary" element={<Summary />} />
        <Route path="/SummeryUser" element={<SummeryUser />} />
      </Routes>
      </BrowserRouter>
    </>
  )
}

export default App
