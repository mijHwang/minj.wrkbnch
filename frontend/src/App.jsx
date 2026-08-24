import { BrowserRouter, Routes, Route } from "react-router-dom";
import Hub from "./pages/Hub";
import CacheDemo from "./pages/CacheDemo";
import ShortenURLDemo from "./pages/ShortenURLDemo";
import "./assets/theme.css";

export default function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Hub />} />
        <Route path="/cache" element={<CacheDemo />} />
        <Route path="/urlShorten" element={<ShortenURLDemo/>}/>
      </Routes>
    </BrowserRouter>
  );
}