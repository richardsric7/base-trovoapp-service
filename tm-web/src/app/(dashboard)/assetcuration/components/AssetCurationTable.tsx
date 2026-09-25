// "use client";
// import CustomTable from "@/components/CustomTable";
// import React, { useCallback, useMemo, useState } from "react";
// import styled from "styled-components";
// import Pagination from "@/components/CustomPagination";
// import curator1 from "@/assets/images/curate1.svg";
// import curator2 from "@/assets/images/curate2.svg";
// import curator3 from "@/assets/images/curate3.svg";
// import editIcon from "@/assets/images/fluent_edit-16-regular.svg";
// import deletIcon from "@/assets/images/deletIcon.svg";
// import positionIcon from "@/assets/images/ion_move.svg";
// import Image from "next/image";
// import CurateAssetModal from "./CurateAssetModal";
// import Link from "next/link";
// import { useRouter } from "next/navigation";
// import DeleteAssetModal from "./DeleteAssetModal";

// const AssetCurationTable = () => {
//   const [currentPage, setCurrentPage] = useState(1);
//   const [pageSize, setPageSize] = useState(10);
// const [openModal, setOpenModal] = useState<boolean>(false);
//   const router = useRouter();
//   const [deleteModal, setDeleteModal] = useState(false);

//   const [moveItem, setMoveItem] = useState(false);
//   // const handleOpenModal = useCallback((record: any) => {
//   //   // setSelectedCuration(record);
//   //   setOpenModal(true);
//   // }, []);

//   // const handleCloseModal = useCallback(() => {
//   //   // setSelectedCuration(null);
//   //   setOpenModal(false);
//   // }, []);
// const columns = [
//   {
//     title: "Position",
//     dataIndex: "position",
//     key: "position",
//   },

//   {
//     title: "Assets",
//     dataIndex: "asset",
//     width: "20%",
//     render: (_: any, record: any) => (
//       <AssetContent>
//         <Image src={curator2} alt="curators" />
//         <div>
//           <AssetName>{record.assetname}</AssetName>
//           <AssetTitle>{record.assettitle}</AssetTitle>
//         </div>
//       </AssetContent>
//     ),
//   },

//   {
//     title: "Created On",
//     dataIndex: "createdOn",
//     key: "createdOn",
//   },

//   {
//     title: "Updated",
//     dataIndex: "updated",
//     key: "updated",
//   },
//   {
//     title: " Curated By",
//     dataIndex: "curatedBy",
//     render: (_: any, record: any) => {
//       return (
//         <CuratorsWrapper>
//           <Image src={curator1} alt="curators" />
//           <CuratorProfile>
//             <CuratorName>{record.curatorName}</CuratorName>
//             <CuratorRole>{record.curatorRole}</CuratorRole>
//           </CuratorProfile>
//         </CuratorsWrapper>
//       );
//     },
//     key: "curatedBy",
//   },

//   {
//     title: "Actions ",
//     dataIndex: "actions",
//     render: (_: any, record: any) => {
//       return (
//         <ActionsWrapper className=" actions-column">
//           <Image
//             src={positionIcon}
//             alt="position-icon"
//             width={16}
//             height={16}
//           />
//           <StyledLink href="/assetcuration/edittoken">
//             <Image src={editIcon} alt="edit" width={16} height={16} />
//           </StyledLink>
//           <Image
//             src={deletIcon}
//             alt="delet-icon"
//             width={16}
//             height={16}
//             onClick={() => setDeleteModal(!deleteModal)}
//           />
//         </ActionsWrapper>
//       );
//     },
//     key: "actions",
//   },
// ];

//   const dataSource = useMemo(() => {
//     return Array.from({ length: 100 }, (_, index) => ({
//       userId: index + 1,
//       position: index + 1,
//       assetname: `AFT (Animal Farm Token)`,
//       assettitle: `Agriculture`,
//       createdOn: `23 Oct 2023`,
//       assestquantity: "12",
//       updated: "2h ago",
//       curatorName: "obi Enechi",
//       curatorRole: "Admin",
//     }));
//   }, []);
//   const paginatedData = useMemo(() => {
//     const startIndex = (currentPage - 1) * pageSize;
//     const endIndex = startIndex + pageSize;
//     return dataSource.slice(startIndex, endIndex);
//   }, [currentPage, pageSize, dataSource]);

//   return (
//     <Container>
//       <CustomTable
//         columns={columns}
//         dataSource={paginatedData}
//         onRowClick={(record) => ({
//           onClick: (event: any) => {
//             // Prevent row click if the user clicks inside the Actions column
//             if (event.target.closest(".actions-column")) return;
//             router.push(`/assetcuration/${record.userId}`);
//           },
//         })}
//       />
//       <Pagination
//         currentPage={currentPage}
//         totalCount={dataSource.length}
//         pageSize={pageSize}
//         onPageChange={setCurrentPage}
//         onPageSizeChange={setPageSize}
//       />
// {openModal && (
//   <CurateAssetModal openModal={openModal} setOpenModal={setOpenModal} />
// )}

// {deleteModal && (
//   <DeleteAssetModal
//     openModal={deleteModal}
//     setOpenModal={setDeleteModal}
//   />
// )}
//     </Container>
//   );
// };

// export default AssetCurationTable;

import React, { useState, useMemo } from "react";
import styled from "styled-components";
import { useRouter } from "next/navigation";
import {
  DndContext,
  closestCenter,
  DragEndEvent,
  DragStartEvent,
} from "@dnd-kit/core";
import {
  SortableContext,
  verticalListSortingStrategy,
  useSortable,
} from "@dnd-kit/sortable";
import { CSS } from "@dnd-kit/utilities";
import Image from "next/image";
import Link from "next/link";
import PositionInputModal from "./PositionInputModal";
import DeleteAssetModal from "./DeleteAssetModal";
import CurateAssetModal from "./CurateAssetModal";
import Pagination from "@/components/CustomPagination";
import curator1 from "@/assets/images/curate1.svg";
import curator2 from "@/assets/images/curate2.svg";
import editIcon from "@/assets/images/fluent_edit-16-regular.svg";
import deletIcon from "@/assets/images/deletIcon.svg";
import positionIcon from "@/assets/images/ion_move.svg";

interface AssetRow {
  id: number;
  position: number;
  assetname: string;
  assettitle: string;
  createdOn: string;
  updated: string;
  curatorName: string;
  curatorRole: string;
}

const SortableRow = ({
  record,
  openPositionModal,
  setDeleteModal,
  isDragging,
  router,
}: {
  record: AssetRow;
  openPositionModal: (record: AssetRow) => void;
  setDeleteModal: (value: boolean) => void;
  isDragging: boolean;
  router: any;
}) => {
  const {
    attributes,
    listeners,
    setNodeRef,
    transform,
    transition,
    isDragging: isRowDragging,
  } = useSortable({ id: record.id });

  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
    backgroundColor: isRowDragging ? "#f5f5f5" : "transparent",
    cursor: "move",
    opacity: isRowDragging ? 0.5 : 1,
  };

  return (
    <Tr
      ref={setNodeRef}
      style={style}
      onClick={() => router.push(`/assetcuration/${record.id}`)}
    >
      <Td {...attributes} {...listeners}>
        {record.position}
      </Td>
      <Td>
        <AssetContent>
          <Image src={curator2} alt="curators" />
          <div>
            <AssetName>{record.assetname}</AssetName>
            <AssetTitle>{record.assettitle}</AssetTitle>
          </div>
        </AssetContent>
      </Td>
      <Td>{record.createdOn}</Td>
      <Td>{record.updated}</Td>
      <Td>
        <CuratorsWrapper>
          <Image src={curator1} alt="curators" />
          <CuratorProfile>
            <CuratorName>{record.curatorName}</CuratorName>
            <CuratorRole>{record.curatorRole}</CuratorRole>
          </CuratorProfile>
        </CuratorsWrapper>
      </Td>
      <Td>
        <ActionsWrapper>
          <PositionIcon
            src={positionIcon}
            alt="position-icon"
            width={16}
            height={16}
            onClick={(e) => {
              e.stopPropagation();
              openPositionModal(record);
            }}
          />
          <StyledLink href="/assetcuration/edittoken">
            <Image src={editIcon} alt="edit" width={16} height={16} />
          </StyledLink>
          <Image
            src={deletIcon}
            alt="delete-icon"
            width={16}
            height={16}
            onClick={() => setDeleteModal(true)}
          />
        </ActionsWrapper>
      </Td>
    </Tr>
  );
};

const AssetCurationTable = () => {
  const [currentPage, setCurrentPage] = useState<number>(1);
  const [pageSize, setPageSize] = useState<number>(10);
  const [deleteModal, setDeleteModal] = useState<boolean>(false);
  const [positionModal, setPositionModal] = useState<boolean>(false);
  const [selectedRow, setSelectedRow] = useState<AssetRow | null>(null);
  const [isDragging, setIsDragging] = useState(false);
  const router = useRouter();
  const [openModal, setOpenModal] = useState<boolean>(false);
  const [dataSource, setDataSource] = useState<AssetRow[]>(
    Array.from({ length: 100 }, (_, index) => ({
      id: index + 1,
      position: index + 1,
      assetname: `AFT (Animal Farm Token)`,
      assettitle: `Agriculture`,
      createdOn: `23 Oct 2023`,
      updated: "2h ago",
      curatorName: "Obi Enechi",
      curatorRole: "Admin",
    }))
  );

  const handleDragStart = (event: DragStartEvent) => {
    setIsDragging(true);
  };

  const handleDragEnd = (event: DragEndEvent) => {
    setIsDragging(false);
    const { active, over } = event;

    if (!active || !over || active.id === over.id) return;

    setDataSource((items) => {
      const oldIndex = items.findIndex((item) => item.id === active.id);
      const newIndex = items.findIndex((item) => item.id === over.id);

      const newItems = [...items];
      const [movedItem] = newItems.splice(oldIndex, 1);
      newItems.splice(newIndex, 0, movedItem);

      return newItems.map((item, index) => ({
        ...item,
        position: index + 1,
      }));
    });
  };

  const openPositionModal = (record: AssetRow) => {
    setSelectedRow(record);
    setPositionModal(true);
  };

  const paginatedData = useMemo(() => {
    const startIndex = (currentPage - 1) * pageSize;
    const endIndex = startIndex + pageSize;
    return dataSource.slice(startIndex, endIndex);
  }, [currentPage, pageSize, dataSource]);

  return (
    <DndContext
      collisionDetection={closestCenter}
      onDragEnd={handleDragEnd}
      onDragStart={handleDragStart}
    >
      <SortableContext
        items={dataSource.map((item) => item.id)}
        strategy={verticalListSortingStrategy}
      >
        <Container>
          <TableWrapper>
            <Table>
              <THead>
                <Tr>
                  <Th>Position</Th>
                  <Th>Assets</Th>
                  <Th>Created On</Th>
                  <Th>Updated</Th>
                  <Th>Curated By</Th>
                  <Th>Actions</Th>
                </Tr>
              </THead>
              <TBody>
                {paginatedData.map((record) => (
                  <SortableRow
                    key={record.id}
                    record={record}
                    openPositionModal={openPositionModal}
                    setDeleteModal={setDeleteModal}
                    isDragging={isDragging}
                    router={router}
                  />
                ))}
              </TBody>
            </Table>
          </TableWrapper>

          <Pagination
            currentPage={currentPage}
            totalCount={dataSource.length}
            pageSize={pageSize}
            onPageChange={setCurrentPage}
            onPageSizeChange={setPageSize}
          />

          {openModal && (
            <CurateAssetModal
              openModal={openModal}
              setOpenModal={setOpenModal}
            />
          )}

          {deleteModal && (
            <DeleteAssetModal
              openModal={deleteModal}
              setOpenModal={setDeleteModal}
            />
          )}

          {positionModal && selectedRow && (
            <PositionInputModal
              isOpen={positionModal}
              onClose={() => setPositionModal(false)}
              maxPosition={dataSource.length}
              setDataSource={setDataSource} //  Pass setDataSource
              selectedRow={selectedRow} //  Pass selectedRow
            />
          )}
        </Container>
      </SortableContext>
    </DndContext>
  );
};

export default AssetCurationTable;

const TableWrapper = styled.div`
  width: 100%;
  overflow-x: auto;
`;

const Table = styled.table`
  width: 100%;
  position: relative;
  th:nth-child(1),
  td:nth-child(1) {
    width: 5%;
  }

  th:nth-child(2) {
    width: 15%;
    text-align: center;
  }
  td:nth-child(2) {
    width: 15%;
  }

  th:nth-child(3),
  td:nth-child(3) {
    width: 10%;
  }

  th:nth-child(4),
  td:nth-child(4) {
    width: 8%;
  }

  th:nth-child(5),
  td:nth-child(5) {
    width: 10%;
  }

  th:nth-child(6),
  td:nth-child(6) {
    width: 4%;
  }
`;

const Tr = styled.tr`
  border-bottom: 1px solid #e5e5ef;
`;

const TBody = styled.tbody``;

const THead = styled.thead`
  border-bottom: 1px solid #e5e5ef;
`;

const Td = styled.td`
  color: #00225a;
  font-size: 12px;
  font-weight: 400;
  line-height: 16px;
  padding: 12px 0;
  border-bottom: 1px solid #e5e5ef;
`;

const Th = styled.th`
  color: #828282;
  font-weight: 500;
  font-size: 14px;
  text-align: left;
  line-height: 16px;
  padding: 12px 0;
  border-bottom: 1px solid #e5e5ef;
`;

const Container = styled.section``;

const CuratorsWrapper = styled.div`
  display: flex;
  align-items: center;
`;

const ActionsWrapper = styled.div`
  display: flex;
  align-items: center;
  gap: 12px;
`;

const AssetName = styled.p`
  font-weight: 500;
  font-size: 14px;
  color: #00225a;
`;

const AssetTitle = styled.p`
  color: #00225a;
  font-weight: 400;
  cursor: pointer;
`;

const StyledLink = styled(Link)`
  text-decoration: none;
  display: block;
`;

const AssetContent = styled.div`
  display: flex;
  align-items: center;
  gap: 2px;
`;

const CuratorProfile = styled.div``;

const CuratorName = styled.p`
  font-size: 14px;
  font-weight: 500;
  line-height: 16px;
  color: #00225a;
`;

const CuratorRole = styled.p`
  font-size: 10px;
  font-weight: 400;
  line-height: 16px;
  color: #00225a;
`;

const DragHandle = styled.div`
  cursor: move;
  display: flex;
  align-items: center;

  &:hover {
    opacity: 0.7;
  }
`;
const PositionIcon = styled(Image)`
  cursor: pointer;
`;
