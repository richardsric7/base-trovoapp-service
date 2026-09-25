import { Modal } from "@/components";
import React from "react";
import styled from "styled-components";
import curator1 from "@/assets/images/curate1.svg";
import curator2 from "@/assets/images/curate2.svg";
import curator3 from "@/assets/images/curate3.svg";
import Image from "next/image";
import assetImg1 from "@/assets/images/realestate.png";
import assetImg2 from "@/assets/images/Ellipse 269.svg";
import editIcon from "@/assets/images/fluent_edit-16-regular.svg";
import deletIcon from "@/assets/images/deletIcon.svg";
import Link from "next/link";

interface CurateModalProps {
  openModal: boolean;
  setOpenModal: React.Dispatch<React.SetStateAction<boolean>>;
  //   onSubmit: () => void;
}

const CurateAssetModal: React.FC<CurateModalProps> = ({
  openModal,
  setOpenModal,
}) => {
  const curatedAssets = [
    {
      id: 1,
      assetImg: assetImg1,
      title: "Atlantis 1",
      description: "Property",
    },
    {
      id: 2,
      assetImg: assetImg2,
      title: "Orchad Estate",
      description: "Property",
    },
    {
      id: 3,
      assetImg: assetImg1,
      title: "Atlantis 1",
      description: "Property",
    },
    {
      id: 4,
      assetImg: assetImg2,
      title: "Orchad Estate",
      description: "Property",
    },
  ];
  return (
    <Modal
      title="First Curation"
      isOpen={openModal}
      onClose={() => setOpenModal(false)}
    >
      <ModalContent>
        <div>
          <AssetCount>
            Total: <HighlightedText>50 Curations</HighlightedText>
          </AssetCount>
          <AssetCount>
            Created On: <HighlightedText>23 oct, 2023</HighlightedText>
          </AssetCount>
          <AssetCounts>
            Curated By:
            <CuratorsWrapper>
              <Image src={curator1} alt="curators" />
              <Image src={curator2} alt="curators" />
              <Image src={curator3} alt="curators" />
            </CuratorsWrapper>
          </AssetCounts>
        </div>

        <ActionsWrapper>
          <StyledLink href="/assetcuration/editcuration">
            <Image src={editIcon} alt="edit" width={16} height={16} />
          </StyledLink>
          <Image src={deletIcon} alt="delet-icon" width={16} height={16} />
        </ActionsWrapper>
      </ModalContent>

      <ModalBody>
        {curatedAssets.map((assets) => (
          <>
            <Wrapper key={assets.id}>
              <Image
                src={assets.assetImg}
                alt="curated-asset"
                width={30}
                height={30}
              />
              <div>
                <Title>{assets.title}</Title>
                <Text>{assets.description}</Text>
              </div>
            </Wrapper>
            <Divider />
          </>
        ))}
      </ModalBody>
    </Modal>
  );
};

export default CurateAssetModal;

const ModalContent = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 18px;
`;

const AssetCount = styled.p`
  font-size: 16px;
  font-weight: 400;
  color: #828282;
  margin: 0;
  line-height: 28px;
`;
const AssetCounts = styled.div`
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 16px;
  font-weight: 400;
  color: #828282;
  margin: 0;
  line-height: 28px;
`;

const HighlightedText = styled.span`
  font-weight: 500;
  color: #00225a;
`;

const CuratorsWrapper = styled.div`
  display: flex;
  align-items: center;
  img {
    margin-left: -8px;
  }

  img:last-child {
    margin-right: 0;
  }
`;

const ActionsWrapper = styled.div`
  display: flex;
  align-items: center;
  gap: 12px;
`;

const ModalBody = styled.div`
  display: flex;
  flex-direction: column;
`;
const Wrapper = styled.div`
  display: flex;
  align-items: center;
  gap: 8px;
`;
const Title = styled.h2`
  font-size: 16px;
  font-weight: 500;
  line-height: 24px;
  letter-spacing: 0.1px;
  color: #00225a;
  margin: 0;
`;
const Text = styled.p`
  font-size: 14px;
  font-weight: 400;
  line-height: 20px;
  letter-spacing: 0.25px;
  color: #00225a;
  margin: 0;
`;
const Divider = styled.div`
  width: 100%;
  height: 1px;
  background-color: #e5e5ef;
  margin: 12px 0;
`;
const StyledLink = styled(Link)`
  text-decoration: none;
  display: block;
`;
