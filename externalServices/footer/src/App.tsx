import React from 'react';
import { FlowText, Link } from './atoms/Texts'
import styled from 'styled-components'
import { EventLink, EventLinkProps} from './atoms/EventButtons'
import { MicropuzzleFrontends } from './config/micro-puzzle-helper';
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
    loading: MicropuzzleFrontends.STARTPAGE_CONTENT,
    text: "Startpage"
  },
  {
    fragmentName: "content",
    loading: MicropuzzleFrontends.ABOUT_CONTENT,
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
