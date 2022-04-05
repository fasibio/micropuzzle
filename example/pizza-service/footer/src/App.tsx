import { useState } from 'react'
import {pushToPage} from './config/micro-puzzle-helper'
import styled from 'styled-components'

const Root = styled.div`
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: space-around;
  background-color: #3a3a3a;
  color: #fafafa;
  `
  const LinkList = styled.ul`
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: flex-start;
    list-style: none;
    padding: 0;
    margin: 0;
  `

const Link = styled.li`
  :hover {
    cursor: pointer;
    text-decoration: underline;
  }
`


function App() {
    return (
      <Root>
        <h1>
          Pizza service micro-puzzle test page
        </h1>
        <LinkList>
          <Link onClick={() => pushToPage("about")}>About</Link>
          <Link onClick={() => pushToPage("about")}>Über</Link>
          <Link onClick={() => pushToPage("about")}>Wer wir sind</Link>
        </LinkList>
      </Root>
  )
}

export default App
