import styled from "styled-components";
import React from "react";
import { LinkButton} from './Buttons'
import {loadMicroFrontend, MicropuzzleFrontends} from '../config/micro-puzzle-helper'
export interface EventLinkProps {
  fragmentName: string
  loading: MicropuzzleFrontends
  text: string
  
}


export const EventLink = ({ fragmentName, loading, text }: EventLinkProps) => {
  return <LinkButton onClick={() => {
    loadMicroFrontend(fragmentName, loading)
  }}>{ text }</LinkButton>
  
}
