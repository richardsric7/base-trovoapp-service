import { useEffect, useState } from 'react';
import {
  TEDropdown,
  TEDropdownToggle,
  TEDropdownMenu,
  TEDropdownItem,
  TERipple,
} from 'tw-elements-react';
import { Asset } from '../types/asset';
import { getAssetCode } from '../utils/utilities';

type AssetDropdownItem = { text: string; value: Asset };

type Props = {
  label: string;
  defaultValue?: AssetDropdownItem | null;
  options: AssetDropdownItem[];
  onSelect: (item: any) => void;
};

export default function AssetDropdown({
  label,
  defaultValue,
  options,
  onSelect,
}: Props) {
  const [selectedItem, setSelectedItem] = useState<AssetDropdownItem | null>(
    defaultValue ?? null,
  );
  useEffect(() => {
    // Update the component when myProp changes
    setSelectedItem(defaultValue ?? null);
  }, [selectedItem]);

  const dropdownItems = options?.map((item: AssetDropdownItem) => (
    <TEDropdownItem
      key={`${new Date().getTime()}${item.text.replace(' ', '')}`}
      id={`${new Date().getTime()}${item.text.replace(' ', '')}`}
      className="bg-primary-200"
    >
      <button
        type="button"
        className="block flex space-x-2 w-full cursor-pointer hover:bg-primary-200 bg-primary-100 whitespace-nowrap px-4 py-2 text-sm text-left font-normal pointer-events-auto active:text-primary-800 focus:hover:bg-primary-200 focus:text-primary-800 focus:outline-none active:no-underline"
        onClick={(e: React.MouseEvent<HTMLButtonElement>) => {
          e.preventDefault;
          onSelect(item);
          setSelectedItem(item);
        }}
      >
        {/* <img src={item.value.imageUrl} alt="logo" /> */}
        <span>{item.text}</span>
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
          {selectedItem?.text ? (
            <div className="flex space-x-2 items-center">
              {/* <img src={selectedItem.value.imageUrl} alt="logo" /> */}
              <span>{selectedItem?.text}</span>
            </div>
          ) : dropdownItems.length ? (
            label
          ) : (
            'No items here'
          )}
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
