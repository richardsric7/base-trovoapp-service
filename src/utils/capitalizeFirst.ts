export default function capitalizeFirstLetter(str: string){
    const s = str?.toLowerCase();
    return s?.charAt(0).toUpperCase() + s?.slice(1);
}