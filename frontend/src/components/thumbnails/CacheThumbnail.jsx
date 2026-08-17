export default function CacheThumbnail() {
  const boxes = [0, 1, 2, 3, 4];
  return (
    <svg viewBox="0 0 300 200" className="cache-thumb-svg">
      {boxes.map((i) => {
        const x = 20 + i * 55;
        const isTail = i === boxes.length - 1;
        return (
          <g key={i} className={isTail ? "node-tail" : "node"}>
            <rect x={x} y="80" width="40" height="40" rx="4" />
            {i < boxes.length - 1 && (
              <line x1={x + 40} y1="100" x2={x + 55} y2="100" markerEnd="url(#arrow)" />
            )}
          </g>
        );
      })}
      <defs>
        <marker id="arrow" markerWidth="6" markerHeight="6" refX="5" refY="3" orient="auto">
          <path d="M0,0 L6,3 L0,6 z" className="arrow-head" />
        </marker>
      </defs>
    </svg>
  );
}