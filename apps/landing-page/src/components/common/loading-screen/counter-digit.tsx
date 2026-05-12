export function CounterDigit({ id, offset }: { id: string; offset?: boolean }) {
  return (
    <div id={id} className={offset ? "offset" : ""}>
      <div className="num">0</div>
      <div className="num">1</div>
      <div className="num">2</div>
      <div className="num">3</div>
      <div className="num">4</div>
      <div className="num">5</div>
      <div className="num">6</div>
      <div className="num">7</div>
      <div className="num">8</div>
      <div className="num">9</div>
    </div>
  );
}
