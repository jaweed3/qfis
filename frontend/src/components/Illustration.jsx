export default function Illustration() {
  return (
    <svg viewBox="0 0 400 400" xmlns="http://www.w3.org/2000/svg" className="w-full h-full">
      {/* Concentric rings */}
      <circle cx="200" cy="200" r="60" fill="none" stroke="white" strokeWidth="0.5" opacity="0.25" />
      <circle cx="200" cy="200" r="100" fill="none" stroke="white" strokeWidth="0.5" opacity="0.2" />
      <circle cx="200" cy="200" r="150" fill="none" stroke="white" strokeWidth="0.5" opacity="0.15" />
      <circle cx="200" cy="200" r="190" fill="none" stroke="white" strokeWidth="0.5" opacity="0.1" />

      {/* Radial lines (24 lines, every 15°) */}
      {Array.from({ length: 24 }).map((_, i) => {
        const angle = (i * 15 * Math.PI) / 180;
        const x2 = 200 + Math.cos(angle) * 190;
        const y2 = 200 + Math.sin(angle) * 190;
        return (
          <line key={i} x1="200" y1="200" x2={x2} y2={y2}
            stroke="white" strokeWidth="0.3" opacity="0.12" />
        );
      })}

      {/* Scales of justice motif — simplified */}
      {/* Beam */}
      <line x1="120" y1="160" x2="280" y2="160" stroke="white" strokeWidth="1.2" opacity="0.6" />
      {/* Center pivot */}
      <line x1="200" y1="160" x2="200" y2="240" stroke="white" strokeWidth="1" opacity="0.5" />
      {/* Left scale */}
      <path d="M120 160 Q120 200 140 200" fill="none" stroke="white" strokeWidth="0.8" opacity="0.5" />
      <path d="M140 200 Q120 200 120 220" fill="none" stroke="white" strokeWidth="0.8" opacity="0.5" />
      {/* Right scale */}
      <path d="M280 160 Q280 200 260 200" fill="none" stroke="white" strokeWidth="0.8" opacity="0.5" />
      <path d="M260 200 Q280 200 280 220" fill="none" stroke="white" strokeWidth="0.8" opacity="0.5" />
      {/* Network nodes on scales */}
      <circle cx="140" cy="210" r="6" fill="none" stroke="white" strokeWidth="0.8" opacity="0.7" />
      <circle cx="260" cy="210" r="6" fill="none" stroke="white" strokeWidth="0.8" opacity="0.7" />
      {/* Node connections */}
      <line x1="140" y1="210" x2="200" y2="240" stroke="white" strokeWidth="0.5" opacity="0.3" strokeDasharray="3,3" />
      <line x1="260" y1="210" x2="200" y2="240" stroke="white" strokeWidth="0.5" opacity="0.3" strokeDasharray="3,3" />

      {/* Small decorative nodes */}
      <circle cx="200" cy="240" r="3" fill="white" opacity="0.4" />
      <circle cx="170" cy="130" r="2.5" fill="none" stroke="white" strokeWidth="0.5" opacity="0.3" />
      <circle cx="240" cy="130" r="2.5" fill="none" stroke="white" strokeWidth="0.5" opacity="0.3" />
      <circle cx="150" cy="280" r="2" fill="none" stroke="white" strokeWidth="0.5" opacity="0.2" />
      <circle cx="250" cy="280" r="2" fill="none" stroke="white" strokeWidth="0.5" opacity="0.2" />

      {/* Crosshatch background lines */}
      {Array.from({ length: 8 }).map((_, i) => (
        <line key={`h1-${i}`}
          x1={i * 50} y1="0" x2={i * 50 + 50} y2="50"
          stroke="white" strokeWidth="0.2" opacity="0.04" />
      ))}
    </svg>
  );
}
