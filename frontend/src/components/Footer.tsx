import Link from "next/link"

export default function Footer() {
    return (
        <footer className="py-10 px-6 md:px-10 bg-gray-100">
            <div className="max-w-6xl mx-auto flex flex-col md:flex-row justify-between items-center">
                <div className="text-2xl font-bold text-blue-600 mb-4 md:mb-0">SwipeRank</div>
                <nav>
                    <ul className="flex space-x-6">
                        <li>
                            <Link href="#" className="text-gray-600 hover:text-blue-600">Privacy Policy</Link>
                        </li>
                        <li>
                            <Link href="#" className="text-gray-600 hover:text-blue-600">Terms of Service</Link>
                        </li>
                        <li>
                            <Link href="#" className="text-gray-600 hover:text-blue-600">Contact Us</Link>
                        </li>
                    </ul>
                </nav>
            </div>
            <div className="mt-6 text-center text-gray-500">© 2025 [Your Project Name]. All rights reserved.</div>
        </footer>
    )
}