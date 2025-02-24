import React from 'react';
import { useNavigate, useLocation } from 'react-router-dom'; // Importa useLocation para obtener los errores

const TablaErrores = () => {
  const navigate = useNavigate();
  const location = useLocation();
  const errores = location.state?.errores || []; // Obtener los errores de la navegación

  return (
    <div>
      <h2 style={{ color: 'white' }}>Tabla de Errores</h2>
      {errores.length > 0 ? (
        <table style={{ width: '100%', color: 'white', backgroundColor: 'rgba(0, 0, 0, 0.5)', borderCollapse: 'collapse' }}>
          <thead>
            <tr>
              <th style={{ border: '1px solid white', color: 'white' }}>#</th> {/* Encabezado para el contador */}
              <th style={{ border: '1px solid white', color: 'white' }}>Tipo de Error</th>
              <th style={{ border: '1px solid white', color: 'white' }}>Descripción</th>
              <th style={{ border: '1px solid white', color: 'white' }}>Línea</th>
              <th style={{ border: '1px solid white', color: 'white' }}>Columna</th>
            </tr>
          </thead>
          <tbody>
            {errores.map((error, index) => (
              <tr key={index}>
                <td style={{ border: '1px solid white', color: 'white' }}>{index + 1}</td> {/* Contador */}
                <td style={{ border: '1px solid white' }}>{error.tipo}</td>
                <td style={{ border: '1px solid white' }}>{error.desc}</td>
                <td style={{ border: '1px solid white' }}>{error.linea}</td>
                <td style={{ border: '1px solid white' }}>{error.col}</td>
              </tr>
            ))}
          </tbody>
        </table>
      ) : (
        <p style={{ color: 'white' }}>No hay errores</p>
      )}

      <button onClick={() => navigate('/')}>Regresar al Menú Principal</button> {/* Botón para regresar */}
    </div>
  );
};

export default TablaErrores;
