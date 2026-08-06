import { useEffect, useState } from 'react';
import {
  TEDropdown,
  TEDropdownToggle,
  TEDropdownMenu,
  TEDropdownItem,
  TERipple,
} from 'tw-elements-react';
import { Wallet } from '../types/wallet';

type WalletDropdownItem = { text: string; index: number; value: Wallet };

type Props = {
  label: string;
  defaultValue?: WalletDropdownItem | null;
  options: WalletDropdownItem[];
  onSelect: (item: any) => void;
};

export default function WalletDropdown({
  label,
  defaultValue,
  options,
  onSelect,
}: Props) {
  const [selectedItem, setSelectedItem] = useState<WalletDropdownItem | null>(
    defaultValue ?? null,
  );
  useEffect(() => {
    // Update the component when myProp changes
    setSelectedItem(defaultValue ?? null);
  }, [selectedItem]);

  const dropdownItems = options?.map((item: WalletDropdownItem) => (
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
        <span>{item.text}</span>
        <span className="flex space-x-1">
          {item.value.sharedAccessEnabled === true && (
            <svg
              xmlns="http://www.w3.org/2000/svg"
              height="20px"
              viewBox="0 -960 960 960"
              width="18px"
              fill="#004988"
            >
              <path d="M96-192v-92q0-25.78 12.5-47.39T143-366q54-32 114.5-49T384-432q66 0 126.5 17T625-366q22 13 34.5 34.61T672-284v92H96Zm648 0v-92q0-42-19.5-78T672-421q39 8 75.5 21.5T817-366q22 13 34.5 34.67Q864-309.65 864-284v92H744ZM384-480q-60 0-102-42t-42-102q0-60 42-102t102-42q60 0 102 42t42 102q0 60-42 102t-102 42Zm336-144q0 60-42 102t-102 42q-8 0-15-.5t-15-2.5q25-29 39.5-64.5T600-624q0-41-14.5-76.5T546-765q8-2 15-2.5t15-.5q60 0 102 42t42 102ZM168-264h432v-20q0-6.47-3.03-11.76-3.02-5.3-7.97-8.24-47-27-99-41.5T384-360q-54 0-106 14t-99 42q-4.95 2.83-7.98 7.91-3.02 5.09-3.02 12V-264Zm216.21-288Q414-552 435-573.21t21-51Q456-654 434.79-675t-51-21Q354-696 333-674.79t-21 51Q312-594 333.21-573t51 21ZM384-264Zm0-360Z" />
            </svg>
          )}
          {item.value.walletType === 1 && (
            <svg
              xmlns="http://www.w3.org/2000/svg"
              height="20px"
              viewBox="0 -960 960 960"
              width="18px"
              fill="#004988"
            >
              <path d="M480-120 156-300v-360l324-180 324 180v360L480-120ZM376-579q20-21 47.5-33t56.5-12q29 0 56.5 12t47.5 33l109-60-213-118-213 118 109 60Zm68 357v-118q-48-12-78-51t-30-89q0-9 .5-18.5T341-516l-113-63v237l216 120Zm36-186q30 0 51-21t21-51q0-30-21-51t-51-21q-30 0-51 21t-21 51q0 30 21 51t51 21Zm36 186 216-120v-237l-113 63q3 8 4 17.5t1 18.5q0 50-30 89t-78 51v118Z" />
            </svg>
          )}
          {item.value.alias.includes('-distribution') && (
            <svg
              xmlns="http://www.w3.org/2000/svg"
              height="20px"
              viewBox="0 -960 960 960"
              width="18px"
              fill="#004988"
            >
              <path d="M624-144v-108H444v-384H336v108H96v-288h240v108h288v-108h240v288H624v-108H516v312h108v-108h240v288H624ZM168-744v144-144Zm528 384v144-144Zm0-384v144-144Zm0 144h96v-144h-96v144Zm0 384h96v-144h-96v144ZM168-600h96v-144h-96v144Z" />
            </svg>
          )}
          {item.value.walletType === 2 && (
            <svg
              xmlns="http://www.w3.org/2000/svg"
              height="20px"
              viewBox="0 -960 960 960"
              width="18px"
              fill="#004988"
            >
              <path d="M264-144q-43 0-75.5-27T146-240h-26q-29.7 0-50.85-21.15Q48-282.3 48-312v-192h432v-144q0-29.7 21.15-50.85Q522.3-720 552-720h72v-72q-2-10 5.47-17 7.48-7 18.53-7h72q10 0 16 7t8 17v72h27q21 0 39 10.5t26 29.5l68 137q4 7.58 6 15.79 2 8.21 2 16.21v271h-98q-10 42-42.5 69T696-144q-43 0-75.5-27T578-240H382q-10 42-42.5 69T264-144Zm0-72q20.4 0 34.2-13.8Q312-243.6 312-264q0-20.4-13.8-34.2Q284.4-312 264-312q-20.4 0-34.2 13.8Q216-284.4 216-264q0 20.4 13.8 34.2Q243.6-216 264-216Zm432 0q20.4 0 34.2-13.8Q744-243.6 744-264q0-20.4-13.8-34.2Q716.4-312 696-312q-20.4 0-34.2 13.8Q648-284.4 648-264q0 20.4 13.8 34.2Q675.6-216 696-216ZM120-432v120h34q14-33 44-52.5t65.5-19.5q35.5 0 65 20t45.7 52H480v-120H120Zm432 120h34q14-33 44-52.5t66-19.5q36 0 66 19.5t44 52.5h34v-120H552v120Zm0-192h288l-69-144H552v144ZM48-552v-48h48v-72H48v-48h384v48h-48v72h48v48H48Zm96-48h72v-72h-72v72Zm120 0h72v-72h-72v72Zm216 168H120h360Zm72 0h288-288Z" />
            </svg>
          )}
        </span>
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
              <span>{selectedItem?.text}</span>
              <span className="flex space-x-1">
                {selectedItem.value.sharedAccessEnabled === true && (
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    height="20px"
                    viewBox="0 -960 960 960"
                    width="18px"
                    fill="#004988"
                  >
                    <path d="M96-192v-92q0-25.78 12.5-47.39T143-366q54-32 114.5-49T384-432q66 0 126.5 17T625-366q22 13 34.5 34.61T672-284v92H96Zm648 0v-92q0-42-19.5-78T672-421q39 8 75.5 21.5T817-366q22 13 34.5 34.67Q864-309.65 864-284v92H744ZM384-480q-60 0-102-42t-42-102q0-60 42-102t102-42q60 0 102 42t42 102q0 60-42 102t-102 42Zm336-144q0 60-42 102t-102 42q-8 0-15-.5t-15-2.5q25-29 39.5-64.5T600-624q0-41-14.5-76.5T546-765q8-2 15-2.5t15-.5q60 0 102 42t42 102ZM168-264h432v-20q0-6.47-3.03-11.76-3.02-5.3-7.97-8.24-47-27-99-41.5T384-360q-54 0-106 14t-99 42q-4.95 2.83-7.98 7.91-3.02 5.09-3.02 12V-264Zm216.21-288Q414-552 435-573.21t21-51Q456-654 434.79-675t-51-21Q354-696 333-674.79t-21 51Q312-594 333.21-573t51 21ZM384-264Zm0-360Z" />
                  </svg>
                )}
                {selectedItem.value.walletType === 1 && (
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    height="20px"
                    viewBox="0 -960 960 960"
                    width="18px"
                    fill="#004988"
                  >
                    <path d="M480-120 156-300v-360l324-180 324 180v360L480-120ZM376-579q20-21 47.5-33t56.5-12q29 0 56.5 12t47.5 33l109-60-213-118-213 118 109 60Zm68 357v-118q-48-12-78-51t-30-89q0-9 .5-18.5T341-516l-113-63v237l216 120Zm36-186q30 0 51-21t21-51q0-30-21-51t-51-21q-30 0-51 21t-21 51q0 30 21 51t51 21Zm36 186 216-120v-237l-113 63q3 8 4 17.5t1 18.5q0 50-30 89t-78 51v118Z" />
                  </svg>
                )}
                {selectedItem.value.alias.includes('-distribution') && (
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    height="20px"
                    viewBox="0 -960 960 960"
                    width="18px"
                    fill="#004988"
                  >
                    <path d="M624-144v-108H444v-384H336v108H96v-288h240v108h288v-108h240v288H624v-108H516v312h108v-108h240v288H624ZM168-744v144-144Zm528 384v144-144Zm0-384v144-144Zm0 144h96v-144h-96v144Zm0 384h96v-144h-96v144ZM168-600h96v-144h-96v144Z" />
                  </svg>
                )}
                {selectedItem.value.walletType === 2 && (
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    height="20px"
                    viewBox="0 -960 960 960"
                    width="18px"
                    fill="#004988"
                  >
                    <path d="M264-144q-43 0-75.5-27T146-240h-26q-29.7 0-50.85-21.15Q48-282.3 48-312v-192h432v-144q0-29.7 21.15-50.85Q522.3-720 552-720h72v-72q-2-10 5.47-17 7.48-7 18.53-7h72q10 0 16 7t8 17v72h27q21 0 39 10.5t26 29.5l68 137q4 7.58 6 15.79 2 8.21 2 16.21v271h-98q-10 42-42.5 69T696-144q-43 0-75.5-27T578-240H382q-10 42-42.5 69T264-144Zm0-72q20.4 0 34.2-13.8Q312-243.6 312-264q0-20.4-13.8-34.2Q284.4-312 264-312q-20.4 0-34.2 13.8Q216-284.4 216-264q0 20.4 13.8 34.2Q243.6-216 264-216Zm432 0q20.4 0 34.2-13.8Q744-243.6 744-264q0-20.4-13.8-34.2Q716.4-312 696-312q-20.4 0-34.2 13.8Q648-284.4 648-264q0 20.4 13.8 34.2Q675.6-216 696-216ZM120-432v120h34q14-33 44-52.5t65.5-19.5q35.5 0 65 20t45.7 52H480v-120H120Zm432 120h34q14-33 44-52.5t66-19.5q36 0 66 19.5t44 52.5h34v-120H552v120Zm0-192h288l-69-144H552v144ZM48-552v-48h48v-72H48v-48h384v48h-48v72h48v48H48Zm96-48h72v-72h-72v72Zm120 0h72v-72h-72v72Zm216 168H120h360Zm72 0h288-288Z" />
                  </svg>
                )}
              </span>
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
