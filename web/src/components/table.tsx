
import { ReactNode } from 'react';
import  styles from './table.module.css';

export interface ColumnDef<T> {
  key: keyof T | string;
  header: string;
  width?: string;
  render?: (row: T) => ReactNode;
  className?: string;
}
// interface PaginationProps {
//   hasNextPage: boolean;
//   totalPages: number;
//   totalCount: number;
//   nextPage: number;
//   hasPreviousPage: boolean;
// }

export interface DataTableProps<T> {
  data: T[];
  columns: ColumnDef<T>[];
  onRowClick?: (row: T) => void;
  onRowAction?: (row: T, action: string) => void;
  itemsPerPageOptions?: number[];
  defaultItemsPerPage?: number;
  showSearch?: boolean;
  searchPlaceholder?: string;
  className?: string;
  emptyMessage?: string;
  onSearch?: (query: string) => void;
  // pagination?: PaginationProps;
  // currentPage?: number;
  // limit?: number;
  // onPageChange?: (page: number) => void;
  // onLimitChange?: (limit: number) => void;
}

export const CustomTable = <T = ReactNode,>({
  data,
  columns,
  onRowClick,
  className = '',
  emptyMessage = 'No data available',
  // pagination,
  // currentPage,
  // limit,
  // onPageChange,
  // onLimitChange,
}: DataTableProps<T>) => {
  return (
<div className={`p-4 bg-[#F7FAFC] rounded-lg ${className}`}>
  <div className="overflow-x-auto">
    <table className="min-w-full text-left text-sm">
      <thead>
        <tr className="text-gray-500 uppercase text-xs border-b">
          {columns.map((column) => (
            <th
              key={column.key as string}
              style={column.width ? { width: column.width } : {}}
              className={`px-4 py-3 font-medium ${column.className || ''}`}
            >
              {column.header}
            </th>
          ))}
        </tr>
      </thead>
      <tbody className={`divide-y divide-gray-100 ${styles.tTody}`}>
        {data.length > 0 ? (
          data.map((row, index) => (
            <tr
              key={index}
              onClick={() => onRowClick && onRowClick(row)}
              className={`bg-white hover:bg-gray-50 transition rounded-lg ${
                onRowClick ? 'cursor-pointer' : ''
              }`}
            >
              {columns.map((column) => (
                <td
                  key={`${index}-${column.key as string}`}
                  className={`px-4 py-4 whitespace-nowrap ${
                    column.className || ''
                  }`}
                >
                  {column.render
                    ? column.render(row)
                    : (row[column.key as keyof T] as ReactNode)}
                </td>
              ))}
            </tr>
          ))
        ) : (
          <tr>
            <td
              colSpan={columns.length}
              className="px-4 py-6 text-center text-gray-400"
            >
              {emptyMessage}
            </td>
          </tr>
        )}
      </tbody>
    </table>
  </div>

  {/* {pagination && (
    <div className="mt-4">
      <CustomPagination
        pagination={pagination}
        onLimitChange={onLimitChange}
        onPageChange={onPageChange}
        limit={limit}
        currentPage={currentPage}
      />
    </div>
  )} */}
</div>
);
};