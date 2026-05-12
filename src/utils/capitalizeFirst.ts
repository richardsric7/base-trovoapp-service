export default function capitalizeFirstLetter(str: string){
    const s = str?.toLowerCase();
    return s?.charAt(0).toUpperCase() + s?.slice(1);
}

export const capitalizeEachWord = (str: string): string => {
  return str
    .toLowerCase()
    .split(' ')
    .map(word => word.charAt(0).toUpperCase() + word.slice(1))
    .join(' ');
};