export default function LinkThumbnail() {
  return (
    <svg viewBox="0 0 300 200" className="cache-thumb-svg">
      <rect x="20" y="85" width="140" height="30" rx="6" className="node" />
      <line x1="170" y1="100" x2="210" y2="100" className="node" markerEnd="url(#linkArrow)" />
      <rect x="220" y="85" width="50" height="30" rx="6" className="node" />
      <defs>
        <marker id="linkArrow" markerWidth="6" markerHeight="6" refX="5" refY="3" orient="auto">
          <path d="M0,0 L6,3 L0,6 z" className="arrow-head" />
        </marker>
      </defs>
    </svg>
  );
}