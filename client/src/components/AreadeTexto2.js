import React from 'react';
import { Editor } from '@monaco-editor/react';

const AreadeTexto2 = ({ responseContent }) => {
    return (
        <Editor
            height="600px"
            width="1400px"
            language="Batch"
            value={responseContent}
            theme="vs-dark"
            options={{ readOnly: true, minimap: { enabled: false } }}
        />
    );
};

export default AreadeTexto2;
