function fn() {
    var movieTitles = [
        'The Midnight Detective', 'Ocean Waves', 'Mountain Sunrise', 'City Lights', 'Desert Storm',
        'Forest Echo', 'River Rapids', 'Sky High', 'Deep Blue', 'Golden Hour',
        'The Last Sunset', 'Echo Valley', 'Silent Night', 'Breaking Dawn', 'The Lost City',
        'Mystery Manor', 'Hidden Truth', 'Final Frontier', 'Beyond the Stars', 'The Great Escape',
        'Shadow Walker', 'Light of Day', 'Dark Waters', 'High Ground', 'Low Tide',
        'The Long Road', 'Short Circuit', 'Fast Forward', 'Slow Motion', 'Time Machine',
        'Space Odyssey', 'Earth Bound', 'Wind Walker', 'Fire Starter', 'Ice Breaker',
        'Storm Chaser', 'Rain Maker', 'Sun Seeker', 'Moon Walker', 'Star Gazer',
        'The Hidden Gem', 'Lost Treasure', 'Ancient Mystery', 'Modern Times', 'Future World',
        'Past Lives', 'Present Tense', 'The Unknown', 'Known Entity', 'Secret Garden',
        'Public Enemy', 'Private Eye', 'Open Door', 'Closed Circuit', 'The Big Picture',
        'Small Details', 'High Stakes', 'Low Profile', 'The Right Stuff', 'Wrong Turn',
        'Left Behind', 'Center Stage', 'The Main Event', 'Side Story', 'Back Story',
        'Front Line', 'The End Game', 'New Beginning', 'Old Friends', 'Young Blood',
        'The First Time', 'Last Chance', 'Next Level', 'Previous Life', 'Current Affairs',
        'The Real Deal', 'Fake News', 'True Story', 'False Alarm', 'The Good Life',
        'Bad Company', 'Ugly Truth', 'Beautiful Mind', 'The Perfect Storm', 'Imperfect World',
        'The Best Years', 'Worst Nightmare', 'Better Days', 'Worse for Wear', 'The Long Goodbye',
        'Short Fuse', 'Tall Order', 'Small World', 'Big Dreams', 'The Little Things',
        'Major League', 'Minor Detail', 'The Grand Tour', 'Petty Crime', 'The Noble Quest',
        'Common Ground', 'Rare Find', 'The Usual Suspects', 'Unusual Circumstances', 'The Normal Life',
        'Abnormal Behavior', 'The Standard Issue', 'Special Delivery', 'The Regular Guy', 'Irregular Hours'
    ];

    var bucketNames = [
        'Action Thrillers', 'Romantic Comedies', 'Sci-Fi Classics', 'Horror Collection', 'Drama Series',
        'Documentary Films', 'Animated Features', 'Foreign Cinema', 'Indie Gems', 'Blockbuster Hits',
        'Cult Classics', 'Award Winners', 'Family Favorites', 'Teen Movies', 'Adult Dramas',
        'Comedy Central', 'Mystery & Suspense', 'Adventure Films', 'War Movies', 'Western Collection',
        'Musical Theater', 'Biographical Films', 'Historical Dramas', 'Fantasy Worlds', 'Crime Stories',
        'Sports Movies', 'Food & Travel', 'Nature Documentaries', 'Space Exploration', 'Ocean Adventures',
        'Mountain Expeditions', 'Urban Stories', 'Rural Life', 'City Nights', 'Country Roads',
        'Beach Vibes', 'Desert Tales', 'Forest Adventures', 'River Journeys', 'Sky Stories', 'Deep Sea',
        'Summer Blockbusters', 'Winter Tales', 'Spring Awakening', 'Autumn Stories', 'Holiday Classics',
        'Christmas Movies', 'Halloween Horror', 'Valentine Romance', 'Easter Stories', 'New Year Films',
        '80s Classics', '90s Nostalgia', '2000s Hits', '2010s Favorites', '2020s New Releases',
        'Golden Age Cinema', 'Silver Screen Classics', 'Modern Masterpieces', 'Contemporary Hits', 'Timeless Tales',
        'Director Spotlight', 'Actor Showcase', 'Genre Collection', 'Theme Nights', 'Mood Movies',
        'Feel Good Films', 'Thought Provoking', 'Mind Bending', 'Heart Warming', 'Soul Searching',
        'Spirit Lifting', 'Mind Expanding', 'Eye Opening', 'Life Changing', 'World Viewing',
        'Cultural Exchange', 'Global Stories', 'Local Tales', 'Regional Cinema', 'National Treasures',
        'International Hits', 'World Cinema', 'Art House Films', 'Experimental Movies', 'Avant Garde',
        'Underground Cinema', 'Independent Films', 'Studio Classics', 'Major Productions', 'Minor Releases',
        'Hidden Gems', 'Forgotten Classics', 'Rediscovered Films', 'Restored Masterpieces', 'Remastered Hits',
        'Director Cuts', 'Extended Editions', 'Special Features', 'Bonus Content', 'Behind the Scenes',
        'Making Of', 'Deleted Scenes', 'Alternate Endings', 'Unrated Versions', 'Rare Footage'
    ];

    
    var genres = [
        'action',     
        'drama',      
        'comedy',     
        'horror',     
        'sci_fi',     
        'romance',     
        'thriller',   
        'fantasy',    
        'documentary', 
        'music',       
        'news',        
        'sports',     
        'kids',       
        'educational'  
    ];

    var getRandomMovieTitle = function() {
        return movieTitles[Math.floor(Math.random() * movieTitles.length)];
    };

    var getRandomBucketName = function() {
        return bucketNames[Math.floor(Math.random() * bucketNames.length)];
    };

    var getRandomGenre = function() {
        return genres[Math.floor(Math.random() * genres.length)];
    };

    var getRandomGenres = function(count) {
        var shuffled = genres.slice().sort(function() { return 0.5 - Math.random(); });
        return shuffled.slice(0, Math.min(count, genres.length));
    };

    var generateAssetSlug = function(title) {
        return title.toLowerCase().replace(/[^a-z0-9]/g, '-').replace(/-+/g, '-').replace(/^-|-$/g, '');
    };

    var generateBucketKey = function(name) {
        return name.toLowerCase().replace(/[^a-z0-9]/g, '-').replace(/-+/g, '-').replace(/^-|-$/g, '');
    };

    return {
        getRandomMovieTitle: getRandomMovieTitle,
        getRandomBucketName: getRandomBucketName,
        getRandomGenre: getRandomGenre,
        getRandomGenres: getRandomGenres,
        generateAssetSlug: generateAssetSlug,
        generateBucketKey: generateBucketKey
    };
} 