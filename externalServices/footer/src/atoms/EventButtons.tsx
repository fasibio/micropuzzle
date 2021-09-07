import styled from "styled-components";
import React from "react";
import { LinkButton} from './Buttons'
export interface EventLinkProps {
  fragmentName: string
  loading: string
  text: string
  
}


export const EventLink = ({ fragmentName, loading, text }: EventLinkProps) => {
  return <LinkButton onClick={() => {
    const event = new CustomEvent('load-content',{
      bubbles: true,
      composed: true,
      cancelable: true,
      detail: {
      loading: loading,
      content: fragmentName
      }
    })
    document.dispatchEvent(event);
  }}>{ text }</LinkButton>
  
}
