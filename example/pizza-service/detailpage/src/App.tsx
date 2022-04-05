import { useState } from 'react'

import styled from 'styled-components'
import { MainDescription } from './MainDescription'

const Root = styled.div`
  height: 90vh; 
  padding: 10px;
  `



function App() {
    return (
    <Root>
      <MainDescription />
      
    </Root>
  )
}

export default App
