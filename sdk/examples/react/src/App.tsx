import React from 'react';
import { ServifyProvider } from '@servify/react';
import ChatDemo from './components/ChatDemo';

const App: React.FC = () => {
  const handleInitialized = () => {
    console.log('Servify SDK initialized');
  };

  const handleError = (error: Error) => {
    console.error('Servify SDK error:', error);
    window.alert(`SDK error: ${error.message}`);
  };

  return (
    <ServifyProvider
      config={{
        apiUrl: 'http://localhost:8080',
        wsUrl: 'ws://localhost:8080/api/v1/ws',
        customerName: 'React User',
        customerEmail: 'react@example.com',
        debug: true,
      }}
      onInitialized={handleInitialized}
      onError={handleError}
    >
      <div className="container">
        <div className="header">
          <h1>Servify</h1>
          <p>React SDK example aligned with the current WebSocket-first contract.</p>
        </div>
        <ChatDemo />
      </div>
    </ServifyProvider>
  );
};

export default App;
