import { Button } from "../components/ui/button"

export default function CTA() {
    return (
        <section className="py-20 px-6 md:px-10 bg-blue-600 text-white text-center">
            <h2 className="text-3xl md:text-4xl font-bold mb-6">Ready to Simplify Your Applicant Ranking?</h2>
            <p className="text-xl mb-10 max-w-2xl mx-auto">Join SwipeRank today and experience a more efficient way to review and rank applicants.</p>
            <Button size="lg" variant="secondary" className="text-blue-600">Sign Up Now</Button>
        </section>
    )
}