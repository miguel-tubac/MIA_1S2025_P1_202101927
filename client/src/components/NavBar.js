import React, { useRef, useState, useEffect } from 'react';
import AreadeTexto1 from './AreadeTexto1';
import AreadeTexto2 from './AreadeTexto2';
//import { useNavigate } from 'react-router-dom';

const NavBar = () => {
  const fileInputRef = useRef(null);
  const [fileContent, setFileContent] = useState('');
  const [responseContent, setResponseContent] = useState('');
  const [errorList, setErrorList] = useState([]);
  //const navigate = useNavigate();

  // Efecto para cargar el estado desde localStorage
  useEffect(() => {
    const savedFileContent = localStorage.getItem('fileContent');
    const savedResponseContent = localStorage.getItem('responseContent');
    const savedErrorList = localStorage.getItem('errorList');

    if (savedFileContent) {
      setFileContent(savedFileContent);
    }
    if (savedResponseContent) {
      setResponseContent(savedResponseContent);
    }
    if (savedErrorList) {
      setErrorList(JSON.parse(savedErrorList));
    }
  }, []);

  // Efecto para guardar el estado en localStorage
  useEffect(() => {
    localStorage.setItem('fileContent', fileContent);
    localStorage.setItem('responseContent', responseContent);
    localStorage.setItem('errorList', JSON.stringify(errorList));
  }, [fileContent, responseContent, errorList]);

  const handleFileButtonClick = () => {
    fileInputRef.current.click();
  };

  const handleFileChange = (event) => {
    const file = event.target.files[0];
    if (file && file.name.endsWith('.smia')) {
      const reader = new FileReader();
      reader.onload = (e) => {
        const content = e.target.result;
        setFileContent(content);
      };
      reader.readAsText(file);
    } else {
      console.log('Por favor selecciona un archivo con extensión .smia');
    }
  };

  const handleExecuteButtonClick = () => {
    sendDataToBackend(fileContent);
  };


  const sendDataToBackend = (data) => {
    const backendUrl = 'http://localhost:4000/interpretar';
    //console.log(data);
    fetch(backendUrl, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ entrada: data }),
    })
      .then((response) => response.json())
      .then((data) => {
        let resultado = data.consola;
        let errores = data.tablaError;
        console.log('Respuesta del backend:', resultado);
        setResponseContent(resultado);
        setErrorList(errores);
      })
      .catch((error) => console.error('Error al enviar datos al backend:', error));
  };

  return (
    <div>
      <nav
        style={{
          marginBottom: '20px',
          backgroundImage: 'linear-gradient(to right, black, purple)',
          padding: '10px',
        }}
      >
        <button style={{ marginRight: '20px' }} onClick={handleFileButtonClick}>
          Archivo
        </button>
        <input
          ref={fileInputRef}
          type="file"
          accept=".smia"
          style={{ display: 'none' }}
          onChange={handleFileChange}
        />
        <button style={{ marginRight: '20px' }} onClick={handleExecuteButtonClick}>
          Ejecutar
        </button>
      </nav>

      <div style={{ display: 'flex', flexDirection: 'column', gap: '20px' }}>
        <AreadeTexto1 fileContent={fileContent} setFileContent={setFileContent} />
        <AreadeTexto2 responseContent={responseContent} />
      </div>
    </div>
  );
};

export default NavBar;
