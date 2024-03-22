import React, { useState } from 'react';
import './App.css'; // Import CSS file for styling
import Output from './Output'

function CodeEditor() {
  const [code, setCode] = useState('');
  const [loggedIn, setLoggedIn] = useState(false);
  const [selectedLanguage, setSelectedLanguage] = useState('javascript');

  const handleLogin = () => {
    setLoggedIn(true);
  };

  const handleLogout = () => {
    setLoggedIn(false);
  };

  const handleChangeLanguage = (language) => {
    setSelectedLanguage(language);
  };

  const handleRunCode = () => {
    // Implement logic to run code
    console.log('Running code:', code);
  };

  const handleProfile = () => {
    // Implement logic to view user profile
    console.log('Viewing user profile');
  };

  const handleSubmitCode = (code, selectedLanguage) => {
    // Define the request payload
    const payload = {
        created_by: "",
        metadata: {
            key1: "value1",
            key2: "value2"
        },
        lang: selectedLanguage,
        solution: code
    };

    // Define the request options
    const requestOptions = {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(payload)
    };

    // Send the POST request
    fetch('http://localhost:8080/submission', requestOptions)
        .then(response => {
            // Handle response
            console.log('Response:', response);
        })
        .catch(error => {
            // Handle error
            console.error('Error:', error);
        });
}


  const handleKeyPress = (e) => {
    if (e.key === 'Enter' && selectedLanguage === 'python') {
      e.preventDefault();
      const cursorPosition = e.target.selectionStart;
      const codeBeforeCursor = code.substring(0, cursorPosition);
      const lastNewLineIndex = codeBeforeCursor.lastIndexOf('\n');
      const previousLine = codeBeforeCursor.substring(lastNewLineIndex + 1); // Previous line
      const indentation = previousLine ? previousLine.match(/^\s*/)[0] : ''; // Indentation from previous line
      const additionalIndentation = previousLine && previousLine.trim().endsWith(':') ? '    ' : ''; // Additional 4 spaces if previous line ends with ':'
      const whitespaceBeforeCursor = (previousLine ? previousLine.substring(0, e.target.selectionStart - lastNewLineIndex - 1) : '').match(/^\s*$/);
      const whitespace = whitespaceBeforeCursor ? whitespaceBeforeCursor[0] : ''; // Extract matched whitespace or set to empty string
      const isNewLineEmpty = previousLine.trim() === '';
      const updatedCode = code.substring(0, cursorPosition) + (isNewLineEmpty ? '\n' : '\n' + indentation + additionalIndentation + whitespace) + code.substring(cursorPosition);
      setCode(updatedCode);
    }
  };

  return (
    <div className="code-editor-container">
      <h1>Code Editor</h1>
      <div className="button-container">
        {!loggedIn ? (
          <button onClick={handleLogin}>Login</button>
        ) : (
          <>
            <button onClick={handleLogout}>Logout</button>
            <button onClick={handleProfile}>Profile</button>
          </>
        )}
        <button onClick={() => handleChangeLanguage('javascript')} className={selectedLanguage === 'javascript' ? 'selected' : ''}>JavaScript</button>
        <button onClick={() => handleChangeLanguage('python')} className={selectedLanguage === 'python' ? 'selected' : ''}>Python</button>
        <button onClick={() => handleChangeLanguage('java')} className={selectedLanguage === 'java' ? 'selected' : ''}>Java</button>
        {/* Add more buttons for other languages */}
        <button onClick={handleRunCode}>Run</button>
        <button onClick={() => handleSubmitCode(code,selectedLanguage)}>Submit Code</button>
      </div>
      <div style={{ overflowX: 'auto' }}> {/* Enable horizontal scrolling */}
        <textarea
          value={code}
          onChange={(e) => setCode(e.target.value)}
          onKeyDown={handleKeyPress}
          rows={23}
          cols={100}
          placeholder="Enter your code here..."
          style={{ whiteSpace: 'pre' }} // Preserve whitespace for proper indentation display
        />
      </div>
      <Output content={"output is here\n yesh\n dshmsd\n"} />
    </div>
  );
}

export default CodeEditor;
