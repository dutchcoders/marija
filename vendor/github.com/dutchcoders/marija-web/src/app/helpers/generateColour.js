import { hash, intToRGB} from './index';

export default function generateColour(str){
    const hashedString = hash(str);
    return `#${intToRGB(hashedString)}`;
}
