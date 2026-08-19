package care.eka.abdm;

import static org.junit.jupiter.api.Assertions.assertEquals;

import java.time.Duration;
import org.junit.jupiter.api.Test;

class ConfigTest {

    @Test
    void baseUrlForProduction() {
        assertEquals("https://api.eka.care", Config.baseUrlFor(Config.PRODUCTION));
    }

    @Test
    void baseUrlForDevelopment() {
        assertEquals("https://api.dev.eka.care", Config.baseUrlFor(Config.DEVELOPMENT));
    }

    @Test
    void baseUrlForUnrecognisedDefaultsToProduction() {
        assertEquals("https://api.eka.care", Config.baseUrlFor("staging"));
    }

    @Test
    void tokenReturnsStaticTokenWhenNoResolverInstalled() throws Exception {
        Config config = new Config(Config.PRODUCTION, "https://api.eka.care", Duration.ofSeconds(1), "test-agent");
        config.setAuthorizationToken("static-token");
        assertEquals("static-token", config.token());
    }

    @Test
    void tokenPrefersResolverOverStaticToken() throws Exception {
        Config config = new Config(Config.PRODUCTION, "https://api.eka.care", Duration.ofSeconds(1), "test-agent");
        config.setAuthorizationToken("static-token");
        config.setTokenResolver(() -> "resolved-token");
        assertEquals("resolved-token", config.token());
    }

    @Test
    void tokenCallsResolverOnEveryCall() throws Exception {
        int[] calls = {0};
        Config config = new Config(Config.PRODUCTION, "https://api.eka.care", Duration.ofSeconds(1), "test-agent");
        config.setTokenResolver(() -> {
            calls[0]++;
            return "token-" + calls[0];
        });
        assertEquals("token-1", config.token());
        assertEquals("token-2", config.token());
        assertEquals(2, calls[0]);
    }
}
