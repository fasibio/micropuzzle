import React, { useState } from "react";
import styled from "styled-components";
import { pushToPage } from "./config/micro-puzzle-helper";
import { Modal } from "./Model";

const Root = styled.div`
display: flex;
flex-direction: row;
gap: 10px;
`

const Image = styled.div`
  height: 450px;
  width: 450px;
  background-color: #3a3a3a;
  display: flex;
  justify-content: center;
  align-items: center;
  color: white;
`

const DescriptionRoot = styled.div`
  display: flex;
  flex: 1;
  flex-direction: column;
  align-items: flex-end;
`

const AmountInput = styled.input`
  text-align: right;
  width: 40px;
`

const AmountRoot = styled.div`
  display: flex;
  gap: 3px;
`

const ModelContent = styled.div`
  display: flex;
  flex-direction: column;
  flex: 1;
`

const ModelHeader = styled.div`
  display: flex;
  flex-direction: column;
  flex: 1;
`;

export const MainDescription = () => {
  const queryParams = new URLSearchParams(window.location.search);
  const id = parseInt(queryParams.get("id")!)
  const [modalOpen, setOpenModal] =useState(false);
  const [amount, setAmount] = useState(1)
  return <Root>
    <Image>Bild</Image>
    <DescriptionRoot>
      <h1>Pizza {id +1 }</h1>
      <p>Preis 7,59€</p>
      <AmountRoot>
        <AmountInput type="number" value={amount} onChange={(v) => setAmount(v.target.value as any)} />
        <button onClick={() => {
          setOpenModal(true)
        }}>Add to cart</button>
      </AmountRoot>
    </DescriptionRoot>
    <Modal open={modalOpen} onClose={() => setOpenModal(false)}>
      <ModelContent>
        <ModelHeader>
        <h1>Thank you for your order!</h1>
        <h1>Pizza in den Warenkorb gelegt</h1>
        </ModelHeader>
        <button onClick={() => {
          setOpenModal(false)
          pushToPage("cart")
        }}>Direkt zum Warenkorb</button>
      </ModelContent>
    </Modal>
  </Root>
}