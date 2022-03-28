
import styled from 'styled-components'
import { loadMicroFrontend, MicropuzzleFrontends } from './config/micro-puzzle-helper'


const Root = styled.div`
  background-color: red;
  height: 50px;
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  align-items: center;
  `

const FlexContainer = styled.div`

`

const CartButton = styled.button`
  background-image: url('./src/assets/cart.svg');
  background-size: 100% 100%;
    width: 50px;
    height: 50px;
`

function App() {
    return (
    <Root>
    <FlexContainer></FlexContainer>      
    <FlexContainer>Pizza...</FlexContainer>
    <FlexContainer>
      <CartButton onClick={() => {
        console.log('hello')
        loadMicroFrontend("content", MicropuzzleFrontends.CART_CONTENT)
      }}/>
    </FlexContainer>      
    </Root>
  )
}

export default App
