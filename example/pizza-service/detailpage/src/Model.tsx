import React from "react";
import styled from "styled-components";

export const ShadowRoot = styled.div`
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
`;

export const ModalRoot = styled.div`
  display: flex;
  flex-direction: column;
  background-color: white;
  border-radius: 10px;
  padding: 10px;
  width: 500px;
  height: 500px;
`;



export type ModalProps ={
  open: boolean
  onClose: () => void
  children: React.ReactNode
}

export const Modal = ({ open, onClose, children }: ModalProps) => {
  if (!open) return null
  return (
    <ShadowRoot onClick={onClose}>
      <ModalRoot>
        {children}
      </ModalRoot>
    </ShadowRoot>
  );
}