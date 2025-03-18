import { useState } from "react";
import { Route, Routes } from "react-router";
import CreateList from "./routes/CreateList";
import styled from "styled-components";
import CreateOrder from "./routes/CreateOrder";

const RoutesContainer = styled.div`
  display: flex;
  flex-direction: column;
  align-items: center;
`;

const App: React.FC = () => (
  <RoutesContainer>
    <Routes>
      <Route index element={<CreateList />} />
      <Route path="/lists/:id" element={<CreateOrder />} />
    </Routes>
  </RoutesContainer>
);

export default App;
