import React, { useEffect, useState } from 'react';
import { Editor } from '@monaco-editor/react';

const AreadeTexto1 = ({ fileContent, setFileContent }) => {
    const [textValue, setTextValue] = useState('');

    useEffect(() => {
        setTextValue(fileContent);
    }, [fileContent]);

    const handleEditorChange = (value) => {
        setTextValue(value);
        setFileContent(value); // Actualizamos el estado fileContent en NavBar también
    };

    return (
        <Editor
            height="600px"
            width="770px"
            language="java"
            value={textValue}
            onChange={handleEditorChange}
            theme="vs-dark" // Puedes cambiar el tema aquí
            options={{ minimap: { enabled: false } }}
        />
    );
};

export default AreadeTexto1;
