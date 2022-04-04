import { useState } from 'react'
import {pushToPage} from './config/micro-puzzle-helper'
function App() {
    return (
      <div>
        <h1>
          Footer
        </h1>
        <button onClick={() => {
          pushToPage('about')
        }}>About</button>
      </div>
  )
}

export default App
