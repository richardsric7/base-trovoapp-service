import { useState } from 'react';
import {
  TEDropdown,
  TEDropdownToggle,
  TEDropdownMenu,
  TEDropdownItem,
  TERipple,
} from 'tw-elements-react';

type Props = {
  label: string;
  options: string[];
  onSelect: (item: any) => void;
};

export default function Dropdown({ label, options, onSelect }: Props) {
  const [selectedItem, setSelectedItem] = useState(null);

  const dropdownItems = options.map((item: any) => (
    <TEDropdownItem
      key={`${new Date().getTime()}${item.replace(' ', '')}`}
      id={`${new Date().getTime()}${item.replace(' ', '')}`}
      className="bg-primary-200"
    >
      <button
        type="button"
        className="block w-full cursor-pointer hover:bg-primary-200 bg-primary-100 whitespace-nowrap px-8 py-2 text-sm text-left font-normal pointer-events-auto active:text-primary-800 focus:hover:bg-primary-200 focus:text-primary-800 focus:outline-none active:no-underline"
        onClick={() => {
          onSelect(item);
          setSelectedItem(item);
        }}
      >
        {item}
      </button>
    </TEDropdownItem>
  ));

  return (
    <TEDropdown className="flex justify-center w-full">
      <TERipple className="w-full" rippleColor="light">
        <TEDropdownToggle className="flex items-center whitespace-nowrap rounded bg-primary-100 hover:bg-primary-200 text-primary-800 px-8 justify-between py-3 rounded-xl w-full">
          {selectedItem ?? label}
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

      <TEDropdownMenu className="w-full bg-primary-100">
        {dropdownItems}
      </TEDropdownMenu>
    </TEDropdown>
  );
}
