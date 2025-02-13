import Link from "next/link";

export default function Footer() {
  return (
    <footer className="fixed bottom-0 left-0 w-full p-4 border-t bg-white">
      <div className="container mx-auto text-center">
        <a
          href="https://github.com/kennedynguyen1/swipe-rank"
          target="_blank"
          rel="noopener noreferrer"
          className="text-gray-600 hover:text-gray-800"
        >
          View on GitHub
        </a>
      </div>
    </footer>
  );
}
