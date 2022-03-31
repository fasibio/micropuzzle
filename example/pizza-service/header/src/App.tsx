
import styled from 'styled-components'
import { pushToPage} from './config/micro-puzzle-helper'


const Root = styled.div`
  background-color: black;
  color: white;
  height: 50px;
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  align-items: center;
  `

const FlexContainer = styled.div`

`

const CartButton = styled.button`
  background-image: url('/assets/cart.svg');
  background-size: contain;
  background-repeat: no-repeat;
  background-position: center;
  background-size: 90% 90%;
  width: 45px;
  height: 45px;
  margin: 2px;
  padding: 10px;
`

function App() {
    return (
    <Root>
    <FlexContainer></FlexContainer>      
    <FlexContainer>Pizza for everyone</FlexContainer>
    <FlexContainer>
      <CartButton title="Cart" onClick={() => {
       pushToPage('cart')
      }}/>
    </FlexContainer>      
    </Root>
  )
}

export default App
