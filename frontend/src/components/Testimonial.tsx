import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"

export default function Testimonial() {
    return (
        <section className="py-20 px-60 md:px-10">
            <Card className="max-w-3xl mx-auto">
                <CardHeader>
                    <CardTitle className="text-2xl font-semibold text-center">What Our Users Say</CardTitle>
                </CardHeader>
                <CardContent className="text-center">
                    <p className="text-lg mb-6">
                        "SwipeRank has revolutionized our applicant review process. It's intuitive, efficient, and has saved us countless hours. Highly recommended for any organization dealing with applications!"
                    </p>
                    <Avatar className="mx-auto mb-2">
                        <AvatarImage src="/placeholder.svg" alt="User" />
                        <AvatarFallback>D</AvatarFallback>
                    </Avatar>
                    <p className="font-semibold">Jane Doe</p>
                    <p className="text-sm text-gray-500">President, Tech Club at University X</p>
                </CardContent>
            </Card>
        </section>
    )
}