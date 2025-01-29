import { Route, BrowserRouter as Router, Routes } from 'react-router-dom';
import Hub from './pages/Hub';
import Login from './pages/Login';
import Profile from './pages/Profile';
import Room from './pages/Room';
//new comment
function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<Login />} />
        <Route path="/hub" element={<Hub />} />
        <Route path="/room/:id" element={<Room />} />
        <Route path="/profile" element={<Profile />} />
      </Routes>
    </Router>
  );
}

export default App;
