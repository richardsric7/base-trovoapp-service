import { useEffect, useState } from 'react';
import {
  TEDropdown,
  TEDropdownToggle,
  TEDropdownMenu,
  TEDropdownItem,
  TERipple,
} from 'tw-elements-react';

type DropdownItem = { text: string; value: any };

type Props = {
  label: string;
  defaultValue?: DropdownItem | null;
  options: DropdownItem[];
  onSelect: (item: any) => void;
};

export default function Dropdown({
  label,
  defaultValue,
  options,
  onSelect,
}: Props) {
  const [selectedItem, setSelectedItem] = useState<DropdownItem | null>(
    defaultValue ?? null,
  );
  useEffect(() => {
    // Keep internal selection in sync when parent default changes.
    setSelectedItem(defaultValue ?? null);
  }, [defaultValue]);

  const dropdownItems = options?.map((item: DropdownItem) => (
    <TEDropdownItem
      key={`${new Date().getTime()}${item.text.replace(' ', '')}`}
      id={`${new Date().getTime()}${item.text.replace(' ', '')}`}
      className="bg-primary-200"
    >
      <button
        type="button"
        className="block w-full cursor-pointer hover:bg-primary-200 bg-primary-100 whitespace-nowrap px-4 py-2 text-sm text-left font-normal pointer-events-auto active:text-primary-800 focus:hover:bg-primary-200 focus:text-primary-800 focus:outline-none active:no-underline"
        onClick={(e: React.MouseEvent<HTMLButtonElement>) => {
          e.preventDefault;
          onSelect(item);
          setSelectedItem(item);
        }}
      >
        {item.text}
      </button>
    </TEDropdownItem>
  ));

  return (
    <TEDropdown className="flex justify-center w-full bg-primary-100">
      <TERipple className="w-full bg-primary-100" rippleColor="light">
        <TEDropdownToggle
          type="button"
          className="flex items-center ring-1 md:ring-2 ring-gray-200 focus-within:ring-primary-600 whitespace-nowrap rounded bg-primary-100 hover:bg-primary-200 text-gray-700 px-5 justify-between py-3 rounded-md h-12 w-full"
        >
          {selectedItem?.text ??
            (dropdownItems.length ? label : 'No items here')}
          <span className="ml-2 [&>svg]:w-5 w-2">
            <svg
              width="10"
              height="6"
              viewBox="0 0 10 6"
              fill="none"
              xmlns="http://www.w3.org/2000/svg"
            >
              <path d="M0 0.289062L5 5.28906L10 0.289062H0Z" fill="#00225A" />
            </svg>
          </span>
        </TEDropdownToggle>
      </TERipple>

      <TEDropdownMenu className="w-full bg-primary-100 max-h-60 overflow-scroll">
        {dropdownItems}
      </TEDropdownMenu>
    </TEDropdown>
  );
}
