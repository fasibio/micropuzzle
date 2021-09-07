import React from 'react';
import { FlowText, Link } from './atoms/Texts'
import styled from 'styled-components'
import { EventLink, EventLinkProps} from './atoms/EventButtons'
const Root = styled.div`
  display: flex;
  align-items:  center;
  justify-content: space-around;
  min-height: 50px;
  background-color: #e3e3e3;
`



const routings: EventLinkProps[] = [
  {
    fragmentName: "content",
    loading: "startpage.content",
    text: "Startpage"
  },
  {
    fragmentName: "content",
    loading: "about.content",
    text: "About"
  },
]

function App() {
  return (
    <Root>
      { routings.map(one => <EventLink {...one} />)}
    </Root>
  );
}

export default App;
