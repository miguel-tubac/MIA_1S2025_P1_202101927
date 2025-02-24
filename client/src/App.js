import React from 'react';
import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import NavBar from './components/NavBar';
import TablaErrores from './components/TablaErrores'; // Importar el componente de TablaErrores

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<NavBar />} />
        <Route path="/errores" element={<TablaErrores />} />
      </Routes>
    </Router>
  );
}

export default App;
