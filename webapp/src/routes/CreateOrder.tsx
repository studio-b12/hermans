import { CLIENT, ShopCategory, ShopData, ShopStoreItem } from "../client";

import Input from "../components/Input";
import RouteContainer from "../components/RouteContainer";
import styled from "styled-components";
import useEffectAsync from "../hooks/useEffectAsync";
import { useParams } from "react-router";
import { useState } from "react";
import { v4 as uuid } from "uuid";

const ProductContainer = styled.div`
  display: flex;
  justify-content: space-between;
  cursor: pointer;
  background-color: #fff3d2;
  padding: 0.2rem 0.4rem;

  &:hover {
    background-color: #fff9e9;
  }

  > div {
    display: flex;
    flex-direction: column;
  }
`;

const CategoryContainer = styled.div`
  display: flex;
  flex-direction: column;
  gap: 0.4rem;

  > h2 {
    margin: 0.5rem 0 0 0.5rem;
    font-size: 1em;
    opacity: 0.75;
    font-weight: normal;
    text-transform: uppercase;
  }
`;

type ProductProps = {
  item: ShopStoreItem;
  onClick: () => void;
};

const Product: React.FC<ProductProps> = ({ item, onClick }) => (
  <ProductContainer onClick={onClick}>
    <div>
      <strong>{item.title}</strong>
      <span>{item.description}</span>
    </div>
    <div>
      <span>{item.price}</span>
    </div>
  </ProductContainer>
);

type CategoryProps = {
  category: ShopCategory;
  onClick: (item: ShopStoreItem) => void;
};

const Category: React.FC<CategoryProps> = ({ category, onClick }) => (
  <CategoryContainer>
    <h2>{category.name}</h2>
    {category.items?.map((i) => (
      <Product key={uuid()} item={i} onClick={() => onClick(i)} />
    ))}
  </CategoryContainer>
);

const CreateOrder: React.FC = () => {
  const { id } = useParams();
  const [shopData, setShopData] = useState<ShopData>();
  const [filter, setFilter] = useState<string>();

  const [selectedFood, setSelectedFood] = useState<ShopStoreItem>();

  useEffectAsync(async () => {
    const shopData = await CLIENT.getShopData();
    setShopData(shopData);
  }, []);

  const filteredShopData = shopData?.categories
    .map((c) => ({
      ...c,
      items: c.items?.filter((i) =>
        i.title.toLowerCase().includes((filter ?? "").toLowerCase())
      ),
    }))
    .filter((c) => c.items?.length ?? 0 > 0);

  return (
    <RouteContainer>
      <h1>Bestellung erstellen</h1>

      <h2>Essen</h2>
      {(selectedFood && (
        <>
          <Product
            item={selectedFood}
            onClick={() => setSelectedFood(undefined)}
          />
        </>
      )) || (
        <>
          <Input
            placeholder="Suche ..."
            value={filter ?? ""}
            onInput={(e) => setFilter(e.currentTarget.value)}
          />
          {filteredShopData &&
            filteredShopData.map((c) => (
              <Category
                key={uuid()}
                category={c}
                onClick={(i) => setSelectedFood(i)}
              />
            ))}
        </>
      )}
    </RouteContainer>
  );
};

export default CreateOrder;
