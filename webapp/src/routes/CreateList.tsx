import { useState } from "react";
import Button from "../components/Button";
import RouteContainer from "../components/RouteContainer";
import { CLIENT, OrderList } from "../client";
import Input from "../components/Input";

const CreateList: React.FC = () => {
  const [creating, setCreating] = useState(false);
  const [list, setList] = useState<OrderList>();

  const createList = async () => {
    setCreating(true);
    const list = await CLIENT.createList();
    setList(list);
    setCreating(false);
  };

  return (
    <RouteContainer>
      <h1>Bestelliste erstellen</h1>
      {list ? (
        <>
          <p>Liste wurde erstellt!</p>
          <Input
            readOnly
            value={`${window.origin}/lists/${list.id}`}
            onFocus={(e) => e.currentTarget.select()}
          />
        </>
      ) : (
        <Button disabled={creating} onClick={createList}>
          {creating ? "Erstelle Liste ..." : "Liste erstellen"}
        </Button>
      )}
    </RouteContainer>
  );
};

export default CreateList;
