import { useState } from 'react'
import styled from 'styled-components'
import { pushToPage } from './config/micro-puzzle-helper'

const Root = styled.div`
  height: 80vh;

`

const Header = styled.div`
  display: flex;
  justify-content: center;
`

const Content = styled.div`
  display: flex;
  align-items: center;
  flex-direction: column;

`
const List = styled.ul`
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  list-style: none;
  padding: 0;
  margin: 0;
  `

const Item = styled.li`
  display: flex;
  justify-content: center;
  align-items: center;
  flex: 1;
  width: 100vh;
  border-top: 1px solid black;

  background-color: #f5f5f5;
  :last-child{
    border-bottom: 1px solid black;
  }
  :hover {
    cursor: pointer;
    background-color: #e5e5e5;
  }
`

function App() {

  const [items, setItems] = useState(['Pizza 1', 'Pizza 2', 'Pizza 3'])

    return (
      <Root>
        <Header>
        <h1>
          We have the best pizza at net 
        </h1>
        </Header>
        <Content>
          <h2>Order your Favorit Pizza:</h2>
          <List>
            {items.map((item, index) => (
              <Item onClick={() => {
                console.log('hier')
                pushToPage(`detail`,{
                  id: `${index}`,
                  title: "blabla"
                })
              }} key={index}>
                <h3>{item}</h3>
              </Item>
            ))}
          </List>
        </Content>
        
    </Root>
  )

}

export default App
