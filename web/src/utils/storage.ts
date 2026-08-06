export const getStorage = (itemKey: string)=>{
    try {
       let value = localStorage.getItem(itemKey);
       if (value) return JSON.parse(value);
     return null
    } catch (error) {
        return null
    }
};
export const setStorage = (itemKey: string, data: Record<string, any>) =>{
    localStorage.setItem(itemKey, JSON.stringify(data));
}
export const clearStorage = () =>{
    localStorage.clear()
}