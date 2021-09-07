import React from 'react';
import { FlowText} from './atoms/Texts'
function App() {
  return (
    <div className="App">
      <header className="App-header">
        <FlowText>
          Edit <code>src/App.tsx</code> and save to reload.
        </FlowText>
        <a
          className="App-link"
          href="https://reactjs.org"
          target="_blank"
          rel="noopener noreferrer"
        >
          Learn React
        </a>
      </header>
    </div>
  );
}

export default App;
